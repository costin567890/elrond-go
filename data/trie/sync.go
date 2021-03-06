package trie

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/storage"
)

var _ data.TrieSyncer = (*trieSyncer)(nil)

type trieNodeInfo struct {
	trieNode node
	received bool
}

type trieSyncer struct {
	rootFound                 bool
	shardId                   uint32
	topic                     string
	rootHash                  []byte
	nodesForTrie              map[string]trieNodeInfo
	waitTimeBetweenRequests   time.Duration
	trie                      *patriciaMerkleTrie
	requestHandler            RequestHandler
	interceptedNodes          storage.Cacher
	mutOperation              sync.RWMutex
	handlerID                 string
	trieSyncStatistics        data.SyncStatisticsHandler
	lastSyncedTrieNode        time.Time
	timeoutBetweenCommits     time.Duration
	maxHardCapForMissingNodes int
}

const maxNewMissingAddedPerTurn = 10
const minTimeoutBetweenNodesCommits = time.Second

// ArgTrieSyncer is the argument for the trie syncer
type ArgTrieSyncer struct {
	RequestHandler                 RequestHandler
	InterceptedNodes               storage.Cacher
	Trie                           data.Trie
	ShardId                        uint32
	Topic                          string
	TrieSyncStatistics             data.SyncStatisticsHandler
	TimeoutBetweenTrieNodesCommits time.Duration
	MaxHardCapForMissingNodes      int
}

// NewTrieSyncer creates a new instance of trieSyncer
func NewTrieSyncer(arg ArgTrieSyncer) (*trieSyncer, error) {
	if check.IfNil(arg.RequestHandler) {
		return nil, ErrNilRequestHandler
	}
	if check.IfNil(arg.InterceptedNodes) {
		return nil, data.ErrNilCacher
	}
	if check.IfNil(arg.Trie) {
		return nil, ErrNilTrie
	}
	if len(arg.Topic) == 0 {
		return nil, ErrInvalidTrieTopic
	}
	if check.IfNil(arg.TrieSyncStatistics) {
		return nil, ErrNilTrieSyncStatistics
	}
	if arg.TimeoutBetweenTrieNodesCommits < minTimeoutBetweenNodesCommits {
		return nil, fmt.Errorf("%w provided: %v, minimum %v",
			ErrInvalidTimeout, arg.TimeoutBetweenTrieNodesCommits, minTimeoutBetweenNodesCommits)
	}
	if arg.MaxHardCapForMissingNodes < 1 {
		return nil, fmt.Errorf("%w provided: %v", ErrInvalidMaxHardCapForMissingNodes, arg.MaxHardCapForMissingNodes)
	}

	pmt, ok := arg.Trie.(*patriciaMerkleTrie)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ts := &trieSyncer{
		requestHandler:            arg.RequestHandler,
		interceptedNodes:          arg.InterceptedNodes,
		trie:                      pmt,
		nodesForTrie:              make(map[string]trieNodeInfo),
		topic:                     arg.Topic,
		shardId:                   arg.ShardId,
		waitTimeBetweenRequests:   time.Second,
		handlerID:                 core.UniqueIdentifier(),
		trieSyncStatistics:        arg.TrieSyncStatistics,
		timeoutBetweenCommits:     arg.TimeoutBetweenTrieNodesCommits,
		maxHardCapForMissingNodes: arg.MaxHardCapForMissingNodes,
	}

	return ts, nil
}

// StartSyncing completes the trie, asking for missing trie nodes on the network
func (ts *trieSyncer) StartSyncing(rootHash []byte, ctx context.Context) error {
	if len(rootHash) == 0 || bytes.Equal(rootHash, EmptyTrieHash) {
		return nil
	}
	if ctx == nil {
		return ErrNilContext
	}

	ts.mutOperation.Lock()
	ts.nodesForTrie = make(map[string]trieNodeInfo)
	ts.nodesForTrie[string(rootHash)] = trieNodeInfo{received: false}
	ts.lastSyncedTrieNode = time.Now()
	ts.mutOperation.Unlock()

	ts.rootFound = false
	ts.rootHash = rootHash

	for {
		shouldRetryAfterRequest, err := ts.checkIfSynced()
		if err != nil {
			return err
		}

		numUnResolved := ts.requestNodes()
		if !shouldRetryAfterRequest && numUnResolved == 0 {
			return nil
		}

		select {
		case <-time.After(ts.waitTimeBetweenRequests):
			continue
		case <-ctx.Done():
			return ErrContextClosing
		}
	}
}

func (ts *trieSyncer) checkIfSynced() (bool, error) {
	var currentNode node
	var err error
	var nextNodes []node
	missingNodes := make(map[string]struct{})
	currentMissingNodes := make([][]byte, 0)
	checkedNodes := make(map[string]struct{})

	newElement := true
	shouldRetryAfterRequest := false

	ts.mutOperation.Lock()
	defer ts.mutOperation.Unlock()

	for newElement {
		newElement = false

		for nodeHash, nodeInfo := range ts.nodesForTrie {
			currentMissingNodes = currentMissingNodes[:0]
			if _, ok := checkedNodes[nodeHash]; ok {
				continue
			}

			currentNode = nodeInfo.trieNode
			if !nodeInfo.received {
				currentNode, err = ts.getNode([]byte(nodeHash))
				if err != nil {
					continue
				}
			}

			checkedNodes[nodeHash] = struct{}{}

			currentMissingNodes, nextNodes, err = currentNode.loadChildren(ts.getNode)
			if err != nil {
				return false, err
			}
			log.Trace("loaded children for node", "hash", currentNode.getHash())

			if len(currentMissingNodes) > 0 {
				hardcapReached := false
				for _, hash := range currentMissingNodes {
					missingNodes[string(hash)] = struct{}{}
					if len(missingNodes) > ts.maxHardCapForMissingNodes {
						hardcapReached = true
						break
					}
				}

				if hardcapReached {
					newElement = false
					break
				}

				nextNodes = append(nextNodes, currentNode)
				_ = ts.addNew(nextNodes)
				shouldRetryAfterRequest = true

				if len(missingNodes) > maxNewMissingAddedPerTurn {
					newElement = false
				}

				continue
			}

			nextNodes, err = currentNode.getChildren(ts.trie.trieStorage.Database())
			if err != nil {
				return false, err
			}

			tmpNewElement := ts.addNew(nextNodes)
			newElement = newElement || tmpNewElement

			delete(ts.nodesForTrie, nodeHash)

			err = encodeNodeAndCommitToDB(currentNode, ts.trie.trieStorage.Database())
			if err != nil {
				return false, err
			}
			ts.resetWatchdog()

			if !ts.rootFound && bytes.Equal([]byte(nodeHash), ts.rootHash) {
				var collapsedRoot node
				collapsedRoot, err = currentNode.getCollapsed()
				if err != nil {
					return false, err
				}

				ts.trie.root = collapsedRoot
				ts.rootFound = true
			}
		}
	}

	if ts.isTimeoutWhileSyncing() {
		return false, ErrTimeIsOut
	}

	for hash := range missingNodes {
		ts.nodesForTrie[hash] = trieNodeInfo{received: false}
	}

	return shouldRetryAfterRequest, nil
}

func (ts *trieSyncer) resetWatchdog() {
	ts.lastSyncedTrieNode = time.Now()
}

func (ts *trieSyncer) isTimeoutWhileSyncing() bool {
	duration := time.Since(ts.lastSyncedTrieNode)
	return duration > ts.timeoutBetweenCommits
}

// adds new elements to needed hash map, lock ts.nodeHashesMutex before calling
func (ts *trieSyncer) addNew(nextNodes []node) bool {
	newElement := false
	for _, nextNode := range nextNodes {
		nextHash := string(nextNode.getHash())

		nodeInfo, ok := ts.nodesForTrie[nextHash]
		if !ok || !nodeInfo.received {
			newElement = true
			ts.trieSyncStatistics.AddNumReceived(1)
			ts.nodesForTrie[nextHash] = trieNodeInfo{
				trieNode: nextNode,
				received: true,
			}
		}
	}

	return newElement
}

// Trie returns the synced trie
func (ts *trieSyncer) Trie() data.Trie {
	return ts.trie
}

func (ts *trieSyncer) getNode(hash []byte) (node, error) {
	nodeInfo, ok := ts.nodesForTrie[string(hash)]
	if ok && nodeInfo.received {
		return nodeInfo.trieNode, nil
	}

	n, ok := ts.interceptedNodes.Get(hash)
	if ok {
		return trieNode(n)
	}

	existingNode, err := getNodeFromDBAndDecode(hash, ts.trie.trieStorage.Database(), ts.trie.marshalizer, ts.trie.hasher)
	if err != nil {
		return nil, ErrNodeNotFound
	}
	err = existingNode.setHash()
	if err != nil {
		return nil, ErrNodeNotFound
	}

	return existingNode, nil
}

func trieNode(data interface{}) (node, error) {
	n, ok := data.(*InterceptedTrieNode)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	return n.node, nil
}

func (ts *trieSyncer) requestNodes() uint32 {
	ts.mutOperation.RLock()
	numUnResolvedNodes := uint32(len(ts.nodesForTrie))
	hashes := make([][]byte, 0)
	for hash, nodeInfo := range ts.nodesForTrie {
		if !nodeInfo.received {
			hashes = append(hashes, []byte(hash))
		}
	}
	ts.trieSyncStatistics.SetNumMissing(ts.rootHash, len(hashes))
	if len(hashes) > 0 {
		ts.requestHandler.RequestTrieNodes(ts.shardId, hashes, ts.topic)
	}
	ts.mutOperation.RUnlock()

	return numUnResolvedNodes
}

// IsInterfaceNil returns true if there is no value under the interface
func (ts *trieSyncer) IsInterfaceNil() bool {
	return ts == nil
}
