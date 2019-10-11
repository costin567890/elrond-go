package shard

import (
	"github.com/ElrondNetwork/elrond-go/core/throttler"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/dataValidators"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/process/factory/containers"
	"github.com/ElrondNetwork/elrond-go/process/rewardTransaction"
	"github.com/ElrondNetwork/elrond-go/process/transaction"
	"github.com/ElrondNetwork/elrond-go/process/unsigned"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

const numGoRoutines = 100

type interceptorsContainerFactory struct {
	accounts              state.AccountsAdapter
	shardCoordinator      sharding.Coordinator
	messenger             process.TopicHandler
	store                 dataRetriever.StorageService
	marshalizer           marshal.Marshalizer
	hasher                hashing.Hasher
	keyGen                crypto.KeyGenerator
	singleSigner          crypto.SingleSigner
	multiSigner           crypto.MultiSigner
	dataPool              dataRetriever.PoolsHolder
	addrConverter         state.AddressConverter
	nodesCoordinator       sharding.NodesCoordinator
	argInterceptorFactory *interceptorFactory.ArgShardInterceptedDataFactory
	globalTxThrottler     process.InterceptorThrottler
	maxTxNonceDeltaAllowed int
	txFeeHandler          process.FeeHandler
	accounts               state.AccountsAdapter
	shardCoordinator       sharding.Coordinator
	messenger              process.TopicHandler
	store                  dataRetriever.StorageService
	marshalizer            marshal.Marshalizer
	hasher                 hashing.Hasher
	keyGen                 crypto.KeyGenerator
	singleSigner           crypto.SingleSigner
	multiSigner            crypto.MultiSigner
	dataPool               dataRetriever.PoolsHolder
	addrConverter          state.AddressConverter
	nodesCoordinator       sharding.NodesCoordinator
	txInterceptorThrottler process.InterceptorThrottler
	maxTxNonceDeltaAllowed int
	txFeeHandler           process.FeeHandler
}

// NewInterceptorsContainerFactory is responsible for creating a new interceptors factory object
func NewInterceptorsContainerFactory(
	accounts state.AccountsAdapter,
	shardCoordinator sharding.Coordinator,
	nodesCoordinator sharding.NodesCoordinator,
	messenger process.TopicHandler,
	store dataRetriever.StorageService,
	marshalizer marshal.Marshalizer,
	hasher hashing.Hasher,
	keyGen crypto.KeyGenerator,
	singleSigner crypto.SingleSigner,
	multiSigner crypto.MultiSigner,
	dataPool dataRetriever.PoolsHolder,
	addrConverter state.AddressConverter,
	maxTxNonceDeltaAllowed int,
	txFeeHandler process.FeeHandler,
) (*interceptorsContainerFactory, error) {
	if accounts == nil || accounts.IsInterfaceNil() {
		return nil, process.ErrNilAccountsAdapter
	}
	if shardCoordinator == nil || shardCoordinator.IsInterfaceNil() {
		return nil, process.ErrNilShardCoordinator
	}
	if messenger == nil || messenger.IsInterfaceNil() {
		return nil, process.ErrNilMessenger
	}
	if store == nil || store.IsInterfaceNil() {
		return nil, process.ErrNilBlockChain
	}
	if marshalizer == nil || marshalizer.IsInterfaceNil() {
		return nil, process.ErrNilMarshalizer
	}
	if hasher == nil || hasher.IsInterfaceNil() {
		return nil, process.ErrNilHasher
	}
	if keyGen == nil || keyGen.IsInterfaceNil() {
		return nil, process.ErrNilKeyGen
	}
	if singleSigner == nil || singleSigner.IsInterfaceNil() {
		return nil, process.ErrNilSingleSigner
	}
	if multiSigner == nil || multiSigner.IsInterfaceNil() {
		return nil, process.ErrNilMultiSigVerifier
	}
	if dataPool == nil || dataPool.IsInterfaceNil() {
		return nil, process.ErrNilDataPoolHolder
	}
	if addrConverter == nil || addrConverter.IsInterfaceNil() {
		return nil, process.ErrNilAddressConverter
	}
	if nodesCoordinator == nil || nodesCoordinator.IsInterfaceNil() {
		return nil, process.ErrNilNodesCoordinator
	}
	if txFeeHandler == nil || txFeeHandler.IsInterfaceNil() {
		return nil, process.ErrNilEconomicsFeeHandler
	}

	txInterceptorThrottler, err := throttler.NewNumGoRoutineThrottler(maxGoRoutineTxInterceptor)
	if err != nil {
		return nil, err
	}

	argInterceptorFactory := &interceptorFactory.ArgShardInterceptedDataFactory{
		ArgMetaInterceptedDataFactory: &interceptorFactory.ArgMetaInterceptedDataFactory{
			Marshalizer:         marshalizer,
			Hasher:              hasher,
			ShardCoordinator:    shardCoordinator,
			MultiSigVerifier:    multiSigner,
			NodesCoordinator:    nodesCoordinator,
		},
		KeyGen:   keyGen,
		Signer:   singleSigner,
		AddrConv: addrConverter,
	}

	icf := &interceptorsContainerFactory{
		accounts:              accounts,
		shardCoordinator:      shardCoordinator,
		messenger:             messenger,
		store:                 store,
		marshalizer:           marshalizer,
		hasher:                hasher,
		keyGen:                keyGen,
		singleSigner:          singleSigner,
		multiSigner:           multiSigner,
		dataPool:              dataPool,
		addrConverter:         addrConverter,
		nodesCoordinator:       nodesCoordinator,
		maxTxNonceDeltaAllowed: maxTxNonceDeltaAllowed,
		argInterceptorFactory: argInterceptorFactory,
		txFeeHandler:           txFeeHandler,
	}

	var err error
	icf.globalTxThrottler, err = throttler.NewNumGoRoutineThrottler(numGoRoutines)
	if err != nil {
		return nil, err
	}

	return icf, nil
}

// Create returns an interceptor container that will hold all interceptors in the system
func (icf *interceptorsContainerFactory) Create() (process.InterceptorsContainer, error) {
	container := containers.NewInterceptorsContainer()

	keys, interceptorSlice, err := icf.generateTxInterceptors()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	keys, interceptorSlice, err = icf.generateUnsignedTxsInterceptors()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	keys, interceptorSlice, err = icf.generateRewardTxInterceptors()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	keys, interceptorSlice, err = icf.generateHdrInterceptor()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	keys, interceptorSlice, err = icf.generateMiniBlocksInterceptors()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	keys, interceptorSlice, err = icf.generateMetachainHeaderInterceptor()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (icf *interceptorsContainerFactory) createTopicAndAssignHandler(
	topic string,
	interceptor process.Interceptor,
	createChannel bool,
) (process.Interceptor, error) {

	err := icf.messenger.CreateTopic(topic, createChannel)
	if err != nil {
		return nil, err
	}

	return interceptor, icf.messenger.RegisterMessageProcessor(topic, interceptor)
}

//------- Tx interceptors

func (icf *interceptorsContainerFactory) generateTxInterceptors() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator

	noOfShards := shardC.NumberOfShards()

	keys := make([]string, noOfShards)
	interceptorSlice := make([]process.Interceptor, noOfShards)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierTx := factory.TransactionTopic + shardC.CommunicationIdentifier(idx)

		interceptor, err := icf.createOneTxInterceptor(identifierTx)
		if err != nil {
			return nil, nil, err
		}

		keys[int(idx)] = identifierTx
		interceptorSlice[int(idx)] = interceptor
	}

	//tx interceptor for metachain topic
	identifierTx := factory.TransactionTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

	interceptor, err := icf.createOneTxInterceptor(identifierTx)
	if err != nil {
		return nil, nil, err
	}

	keys = append(keys, identifierTx)
	interceptorSlice = append(interceptorSlice, interceptor)
	return keys, interceptorSlice, nil
}

func (icf *interceptorsContainerFactory) createOneTxInterceptor(identifier string) (process.Interceptor, error) {
	txValidator, err := dataValidators.NewTxValidator(icf.accounts, icf.shardCoordinator, icf.maxTxNonceDeltaAllowed)
func (icf *interceptorsContainerFactory) createOneTxInterceptor(topic string) (process.Interceptor, error) {
	txValidator, err := dataValidators.NewTxValidator(icf.accounts, icf.shardCoordinator)
	if err != nil {
		return nil, err
	}

	argProcessor := &processor.ArgTxInterceptorProcessor{
		ShardedDataCache: icf.dataPool.Transactions(),
		TxValidator:      txValidator,
	}
	txProcessor, err := processor.NewTxInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, err
	}

	txFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedTx,
	)
	if err != nil {
		return nil, err
	}

	interceptor, err := interceptors.NewMultiDataInterceptor(
		icf.marshalizer,
		txFactory,
		txProcessor,
		icf.globalTxThrottler,
	)
	if err != nil {
		return nil, err
	}

	return icf.createTopicAndAssignHandler(topic, interceptor, true)
}

	//------- Reward transactions interceptors

	func (icf *interceptorsContainerFactory) generateRewardTxInterceptors() ([]string, []process.Interceptor, error) {
		shardC := icf.shardCoordinator

		noOfShards := shardC.NumberOfShards()

		keys := make([]string, noOfShards)
		interceptorSlice := make([]process.Interceptor, noOfShards)

		for idx := uint32(0); idx < noOfShards; idx++ {
			identifierScr := factory.RewardsTransactionTopic + shardC.CommunicationIdentifier(idx)

			interceptor, err := icf.createOneRewardTxInterceptor(identifierScr)
			if err != nil {
				return nil, nil, err
			}

			keys[int(idx)] = identifierScr
			interceptorSlice[int(idx)] = interceptor
		}

		identifierTx := factory.RewardsTransactionTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

		interceptor, err := icf.createOneRewardTxInterceptor(identifierTx)
		if err != nil {
			return nil, nil, err
		}

		keys = append(keys, identifierTx)
		interceptorSlice = append(interceptorSlice, interceptor)

		return keys, interceptorSlice, nil
	}

	func (icf *interceptorsContainerFactory) createOneRewardTxInterceptor(identifier string) (process.Interceptor, error) {
		rewardTxStorer := icf.store.GetStorer(dataRetriever.RewardTransactionUnit)

		interceptor, err := rewardTransaction.NewRewardTxInterceptor(
			icf.marshalizer,
			icf.dataPool.RewardTransactions(),
			rewardTxStorer,
			icf.addrConverter,
			icf.hasher,
			icf.shardCoordinator,
		)


		//------- Unsigned transactions interceptors

func (icf *interceptorsContainerFactory) generateUnsignedTxsInterceptors() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator

	noOfShards := shardC.NumberOfShards()

	keys := make([]string, noOfShards)
	interceptorSlice := make([]process.Interceptor, noOfShards)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierScr := factory.UnsignedTransactionTopic + shardC.CommunicationIdentifier(idx)

		interceptor, err := icf.createOneUnsignedTxInterceptor(identifierScr)
		if err != nil {
			return nil, nil, err
		}

		keys[int(idx)] = identifierScr
		interceptorSlice[int(idx)] = interceptor
	}

	identifierTx := factory.UnsignedTransactionTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

	interceptor, err := icf.createOneUnsignedTxInterceptor(identifierTx)
	if err != nil {
		return nil, nil, err
	}

	keys = append(keys, identifierTx)
	interceptorSlice = append(interceptorSlice, interceptor)
	return keys, interceptorSlice, nil
}

func (icf *interceptorsContainerFactory) createOneUnsignedTxInterceptor(topic string) (process.Interceptor, error) {
	//TODO replace the nil tx validator with white list validator
	txValidator, err := mock.NewNilTxValidator()
	if err != nil {
		return nil, err
	}

	argProcessor := &processor.ArgTxInterceptorProcessor{
		ShardedDataCache: icf.dataPool.UnsignedTransactions(),
		TxValidator:      txValidator,
	}
	txProcessor, err := processor.NewTxInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, err
	}

	txFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedUnsignedTx,
	)
	if err != nil {
		return nil, err
	}

	interceptor, err := interceptors.NewMultiDataInterceptor(
		icf.marshalizer,
		txFactory,
		txProcessor,
		icf.globalTxThrottler,
	)
	if err != nil {
		return nil, err
	}

	return icf.createTopicAndAssignHandler(topic, interceptor, true)
}

//------- Hdr interceptor

func (icf *interceptorsContainerFactory) generateHdrInterceptor() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator
	//TODO implement other HeaderHandlerProcessValidator that will check the header's nonce
	// against blockchain's latest nonce - k finality
	hdrValidator, err := dataValidators.NewNilHeaderValidator()
	if err != nil {
		return nil, nil, err
	}

	hdrFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedShardHeader,
	)

	argProcessor := &processor.ArgHdrInterceptorProcessor{
		Headers:       icf.dataPool.Headers(),
		HeadersNonces: icf.dataPool.HeadersNonces(),
		HdrValidator:  hdrValidator,
	}
	hdrProcessor, err := processor.NewHdrInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, nil, err
	}

	//only one intrashard header topic
	interceptor, err := interceptors.NewSingleDataInterceptor(
		hdrFactory,
		hdrProcessor,
		icf.globalTxThrottler,
	)
	if err != nil {
		return nil, nil, err
	}

	identifierHdr := factory.HeadersTopic + shardC.CommunicationIdentifier(shardC.SelfId())
	_, err = icf.createTopicAndAssignHandler(identifierHdr, interceptor, true)
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierHdr}, []process.Interceptor{interceptor}, nil
}

//------- MiniBlocks interceptors

func (icf *interceptorsContainerFactory) generateMiniBlocksInterceptors() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator
	noOfShards := shardC.NumberOfShards()
	keys := make([]string, noOfShards)
	interceptorSlice := make([]process.Interceptor, noOfShards)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(idx)

		interceptor, err := icf.createOneMiniBlocksInterceptor(identifierMiniBlocks)
		if err != nil {
			return nil, nil, err
		}

		keys[int(idx)] = identifierMiniBlocks
		interceptorSlice[int(idx)] = interceptor
	}

	return keys, interceptorSlice, nil
}

func (icf *interceptorsContainerFactory) createOneMiniBlocksInterceptor(topic string) (process.Interceptor, error) {
	argProcessor := &processor.ArgTxBodyInterceptorProcessor{
		MiniblockCache:   icf.dataPool.MiniBlocks(),
		Marshalizer:      icf.marshalizer,
		Hasher:           icf.hasher,
		ShardCoordinator: icf.shardCoordinator,
	}
	txBlockBodyProcessor, err := processor.NewTxBodyInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, err
	}

	txFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedTxBlockBody,
	)
	if err != nil {
		return nil, err
	}

	interceptor, err := interceptors.NewSingleDataInterceptor(
		txFactory,
		txBlockBodyProcessor,
		icf.globalTxThrottler,
	)
	if err != nil {
		return nil, err
	}

	return icf.createTopicAndAssignHandler(topic, interceptor, true)
}

//------- MetachainHeader interceptors

func (icf *interceptorsContainerFactory) generateMetachainHeaderInterceptor() ([]string, []process.Interceptor, error) {
	identifierHdr := factory.MetachainBlocksTopic
	//TODO implement other HeaderHandlerProcessValidator that will check the header's nonce
	// against blockchain's latest nonce - k finality
	hdrValidator, err := dataValidators.NewNilHeaderValidator()
	if err != nil {
		return nil, nil, err
	}

	hdrFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedMetaHeader,
	)

	argProcessor := &processor.ArgHdrInterceptorProcessor{
		Headers:       icf.dataPool.MetaBlocks(),
		HeadersNonces: icf.dataPool.HeadersNonces(),
		HdrValidator:  hdrValidator,
	}
	hdrProcessor, err := processor.NewHdrInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, nil, err
	}

	//only one metachain header topic
	interceptor, err := interceptors.NewSingleDataInterceptor(
		hdrFactory,
		hdrProcessor,
		icf.globalTxThrottler,
	)
	if err != nil {
		return nil, nil, err
	}

	_, err = icf.createTopicAndAssignHandler(identifierHdr, interceptor, true)
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierHdr}, []process.Interceptor{interceptor}, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (icf *interceptorsContainerFactory) IsInterfaceNil() bool {
	if icf == nil {
		return true
	}
	return false
}
