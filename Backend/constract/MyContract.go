// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DAppMeetingParticipant is an auto generated low-level Go binding around an user-defined struct.
type DAppMeetingParticipant struct {
	WalletAddress common.Address
	Name          string
	SessionID     string
}

// DAppMeetingTrack is an auto generated low-level Go binding around an user-defined struct.
type DAppMeetingTrack struct {
	TrackName   string
	Mid         string
	Location    string
	IsPublished bool
	SessionId   string
	RoomId      string
}

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"eventData\",\"type\":\"bytes\"}],\"name\":\"EventForwardedToBackend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"eventData\",\"type\":\"bytes\"}],\"name\":\"EventForwardedToFrontend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"trackName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"mid\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"location\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isPublished\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"sessionId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"}],\"indexed\":false,\"internalType\":\"structDAppMeeting.Track[]\",\"name\":\"initialTracks\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"sessionDescription\",\"type\":\"bytes\"}],\"name\":\"ParticipantJoined\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"}],\"name\":\"ParticipantLeft\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"trackName\",\"type\":\"string\"}],\"name\":\"TrackAdded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_backend\",\"type\":\"address\"}],\"name\":\"addAuthorizedBackend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"trackName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"mid\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"location\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isPublished\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"sessionId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"}],\"internalType\":\"structDAppMeeting.Track\",\"name\":\"_newTrack\",\"type\":\"tuple\"}],\"name\":\"addTrack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"authorizedBackends\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"}],\"name\":\"createRoom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"_eventData\",\"type\":\"bytes\"}],\"name\":\"forwardEventToBackend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_participant\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_eventData\",\"type\":\"bytes\"}],\"name\":\"forwardEventToFrontend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"}],\"name\":\"getParticipantInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"walletAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"sessionID\",\"type\":\"string\"}],\"internalType\":\"structDAppMeeting.Participant\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_participant\",\"type\":\"address\"}],\"name\":\"getParticipantTracks\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"trackName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"mid\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"location\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isPublished\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"sessionId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"}],\"internalType\":\"structDAppMeeting.Track[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_participant\",\"type\":\"address\"}],\"name\":\"getParticipantTracksCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"}],\"name\":\"getRoomParticipantsCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"trackName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"mid\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"location\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isPublished\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"sessionId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"}],\"internalType\":\"structDAppMeeting.Track[]\",\"name\":\"_initialTracks\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"sessionDescription\",\"type\":\"bytes\"}],\"name\":\"joinRoom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"}],\"name\":\"leaveRoom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"participantIndices\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"participantTrackCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"participantTracks\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"trackName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"mid\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"location\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isPublished\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"sessionId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"participantsInRoom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_backend\",\"type\":\"address\"}],\"name\":\"removeAuthorizedBackend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"rooms\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"roomId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"creationTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_roomId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_participantAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_sessionID\",\"type\":\"string\"}],\"name\":\"setParticipantSessionID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080806040523460a657600580546001600160a01b03191633179055600654680100000000000000008110156092576001810180600655811015607e5760065f527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0180546001600160a01b031916331790556123ce90816100ab8239f35b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52604160045260245ffd5b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c806262b7481461167b5780632aba9dda1461162e5780633b4ddda8146114ea5780633c229779146111705780633f2cc59f146110f757806347100107146110b55780635b85facc14610f915780637306d2dd14610dd35780638da5cb5b14610dab5780638df5894514610d635780639ba78fd914610ca5578063aaea818014610bc7578063bb9af00d146109cb578063bbbcc86914610952578063c3241b95146107fe578063e970cf2c146106bf578063f2776d2e14610649578063f3398cf61461051e578063f7cc8724146104835763fab38543146100f4575f80fd5b3461047f57602036600319011261047f576004356001600160401b03811161047f5761012490369060040161171a565b60405161015161014a835192602081818701958087835e81015f81520301902054611914565b1515611ad3565b60405160208184518085835e8101600281520301902060018060a01b0333165f5260205261018560ff60405f205416611b15565b60405160208184518085835e8101600181520301902060018060a01b0333165f5260205260405f2054600260405160208186518087835e81015f81520301902001545f19810190811161038f5781106103ca575b50600260405160208185518086835e81015f81520301902001805480156103b6575f19019061020882826120bb565b6103a3576002815f61022793556102216001820161234a565b0161234a565b5560405160208184518085835e8101600281520301902060018060a01b0333165f5260205260405f2060ff19815416905560405160208184518085835e8101600181520301902060018060a01b0333165f526020525f604081205560405160208184518085835e8101600381520301902060018060a01b0333165f5260205260405f208054905f815581610326575b7fb4d5a3866c2b36a076b8cac89e6068deeff68398781315f07060bdead654e63b61031b85856020604051809284518091835e8101600481520301902060018060a01b0333165f526020525f60408120556040519182916040835260408301906119ec565b3360208301520390a1005b8160060291600683040361038f575f5260205f20908101905b818110156102b6578061035360069261234a565b61035f6001820161234a565b61036b6002820161234a565b5f600382015561037d6004820161234a565b6103896005820161234a565b0161033f565b634e487b7160e01b5f52601160045260245ffd5b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52603160045260245ffd5b600260405160208186518087835e81015f81520301902001600260405160208187518088835e81015f8152030190200154905f19820191821161038f5761041a91610414916120bb565b50612311565b6104458161043f8460026040516020818b51808c835e81015f815203019020016120bb565b906120d4565b60405160208186518087835e60019082019081520301902090516001600160a01b03165f90815260209190915260408120919091556101d9565b5f80fd5b3461047f5761049136611799565b906020604051916104ba61014a8251948481818601978089835e81015f81520301902054611914565b604051828183518087835e8101600281520301902060018060a01b0385165f5282526104ec60ff60405f2054166122c5565b604051928391518091835e810160048152030190209060018060a01b03165f52602052602060405f2054604051908152f35b3461047f57606036600319011261047f576004356001600160401b03811161047f5761054e90369060040161171a565b61055661176d565b6044356001600160401b03811161047f5761057590369060040161171a565b9160405161059c61014a835192602081818701958087835e81015f81520301902054611914565b6020604051809284518091835e8101600281520301902060018060a01b0383165f5260205260ff60405f20541615610604576105ff7fb588e6e55cc56c9ac78812d33d4f8fbbebe720ee569bed8fd4615ddb178fcb129360405193849384612084565b0390a1005b60405162461bcd60e51b815260206004820152601e60248201527f546172676574207061727469636970616e74206e6f7420696e20726f6f6d00006044820152606490fd5b3461047f57602036600319011261047f57610662611783565b60065490600160401b8210156106ab576106858260016106a994016006556118cf565b81546001600160a01b0393841660039290921b91821b9390911b1916919091179055565b005b634e487b7160e01b5f52604160045260245ffd5b3461047f57602036600319011261047f576004356001600160401b03811161047f576107b36104146106f86107fa93369060040161171a565b606060408051610707816116c3565b5f81528260208201520152600260405161073a61014a845192602081818801958087835e81015f81520301902054611914565b60405160208185518085835e81018581520301902060018060a01b0333165f5260205261076d60ff60405f205416611b15565b60405160208185518085835e8101600181520301902060018060a01b0333165f52602052602060405f205493604051928391518091835e81015f815203019020016120bb565b6040519182916020835260018060a01b03815116602084015260406107e6602083015160608387015260808601906119ec565b910151838203601f190160608501526119ec565b0390f35b3461047f5761080c36611799565b9060206040519161083561014a8251948481818601978089835e81015f81520301902054611914565b604051828183518087835e8101600281520301902060018060a01b0385165f52825261086760ff60405f2054166122c5565b604051928391518091835e810160038152030190209060018060a01b03165f5260205260405f208054610899816118b8565b916108a760405193846116f9565b81835260208301905f5260205f205f915b8383106108d557604051602080825281906107fa90820188611a10565b600660206001926040516108e8816116de565b6108f18661194c565b81526108fe85870161194c565b8382015261090e6002870161194c565b604082015260ff6003870154161515606082015261092e6004870161194c565b608082015261093f6005870161194c565b60a08201528152019201920191906108b8565b3461047f57602036600319011261047f576004356001600160401b03811161047f576020806109886109c193369060040161171a565b604051928184925191829101835e81015f81520301902060016109aa8261194c565b9101546040519283926040845260408401906119ec565b9060208301520390f35b3461047f57606036600319011261047f576004356001600160401b03811161047f576109fb90369060040161171a565b610a0361176d565b906044356001600160401b03811161047f5760029182610a2a610acb93369060040161171a565b9460405190610a5261014a855193602081818901968088835e81015f81520301902054611914565b60405160208186518086835e81018681520301902060018060a01b0382165f52602052610a8560ff60405f2054166122c5565b60405160208186518086835e810160018152030190209060018060a01b03165f52602052602060405f205493604051928391518091835e81015f815203019020016120bb565b500181516001600160401b0381116106ab57610ae78254611914565b601f8111610b8c575b50602092601f8211600114610b3057610b21929382915f92610b25575b50508160011b915f199060031b1c19161790565b9055005b015190508480610b0d565b601f19821693835f52805f20915f5b868110610b745750836001959610610b5c575b505050811b019055005b01515f1960f88460031b161c19169055838080610b52565b91926020600181928685015181550194019201610b3f565b610bb790835f5260205f20601f840160051c81019160208510610bbd575b601f0160051c0190611b61565b83610af0565b9091508190610baa565b3461047f57604036600319011261047f576004356001600160401b03811161047f57610bf790369060040161171a565b602435906001600160401b03821161047f57610c387f5af2a6a9c8117113b24e80f03488d6bb6d93973b5faed013eca8657949140f1e92369060040161171a565b90604051610c5f61014a835192602081818701958087835e81015f81520301902054611914565b6020604051809284518091835e8101600281520301902060018060a01b0333165f52602052610c9460ff60405f205416611b15565b6105ff604051928392339084612084565b3461047f57602036600319011261047f57610cbe611783565b6001600160a01b03165f5b600654808210156106a95782610cde836118cf565b905460039190911b1c6001600160a01b031614610cfe5750600101610cc9565b5f1981019250821161038f57610685610d19610d31936118cf565b905460039190911b1c6001600160a01b0316916118cf565b60065480156103b6575f1901610d46816118cf565b81546001600160a01b0360039290921b9190911b19169055600655005b3461047f57602080610d7436611799565b9290604051928184925191829101835e810160018152030190209060018060a01b03165f52602052602060405f2054604051908152f35b3461047f575f36600319011261047f576005546040516001600160a01b039091168152602090f35b3461047f57602036600319011261047f576004356001600160401b03811161047f57610e0390369060040161171a565b60405190610e27815192602081818501958087835e81015f81520301902054611914565b610f565760405160208183518086835e81015f815203019020918151926001600160401b0384116106ab57610e5c8154611914565b601f8111610f26575b50602093601f8111600114610ec35780610e9791600195965f91610eb8575b508160011b915f199060031b1c19161790565b90555b60405192518091845e8201915f835260208142940301902001555f80f35b905086015187610e84565b93601f19851690825f52805f20915f5b818110610f0e575091869160019697879410610ef6575b5050811b019055610e9a565b8701515f1960f88460031b161c191690558680610eea565b86830151845560019093019260209283019201610ed3565b610f5090825f5260205f20601f870160051c81019160208810610bbd57601f0160051c0190611b61565b84610e65565b60405162461bcd60e51b8152602060048201526013602482015272526f6f6d20616c72656164792065786973747360681b6044820152606490fd5b3461047f57606036600319011261047f576004356001600160401b03811161047f57610fc190369060040161171a565b610fc961176d565b604051825160443593602091839181908401835e810160038152030190209060018060a01b03165f5260205260405f20805482101561047f576110769161100f916118fb565b506107fa61101c8261194c565b916110a761102c6001830161194c565b916110396002820161194c565b9061109260ff60038301541692611084611061600561105a6004870161194c565b950161194c565b966040519a8b9a60c08c5260c08c01906119ec565b908a820360208c01526119ec565b9088820360408a01526119ec565b911515606087015285820360808701526119ec565b9083820360a08501526119ec565b3461047f57602036600319011261047f5760043560065481101561047f576110de6020916118cf565b905460405160039290921b1c6001600160a01b03168152f35b3461047f57602036600319011261047f576004356001600160401b03811161047f57600261112b602092369060040161171a565b826040519161115261014a8251948481818601978089835e81015f81520301902054611914565b604051928391518091835e81015f8152030190200154604051908152f35b3461047f57608036600319011261047f576004356001600160401b03811161047f576111a090369060040161171a565b6024356001600160401b03811161047f576111bf90369060040161171a565b90604435906001600160401b03821161047f573660238301121561047f5781600401356111eb816118b8565b926111f960405194856116f9565b8184526024602085019260051b8201019036821161047f5760248101925b8284106114bb57505050506064356001600160401b03811161047f5761124190369060040161171a565b6040519361126861014a84519660208181880199808b835e81015f81520301902054611914565b60405160208185518089835e8101600281520301902060018060a01b0333165f5260205260ff60405f20541661148457604051906112a5826116c3565b33825260208201526020906040516112bd83826116f9565b5f81526040820152600260405183818751808b835e81015f8152030190200190815491600160401b8310156106ab578261043f9160016112ff950181556120bb565b600260405182818651808a835e81015f81520301902001545f19810190811161038f5760409492945182818651808a835e8101600181520301902060018060a01b0333165f52825260405f2055604051818185518089835e600290820190815203019020335f81815291835260408220805460ff1916600117905592905b8151811015611421576001907f989240cc5cb309e91c68e61536b623e192afd2de138e3f69b3eed730d574bfda611418858360051b8601015160405187818d8c518091835e81016003815203019020885f5287526113de8160405f20611b77565b60405187818d8c518091835e81016004815203019020885f52875260405f206114078154612076565b905551604051918291338b84612084565b0390a10161137d565b7f21e27f169d16f9bc6b07bf8a2c343f07f22bfd14122a3181cf25419574f0109b611463866105ff8961147688886040519687966080885260808801906119ec565b9133908701528582036040870152611a10565b9083820360608501526119ec565b60405162461bcd60e51b815260206004820152600f60248201526e416c726561647920696e20726f6f6d60881b6044820152606490fd5b83356001600160401b03811161047f576020916114df8392602436918701016117da565b815201930192611217565b3461047f57604036600319011261047f576004356001600160401b03811161047f5761151a90369060040161171a565b602435906001600160401b03821161047f5761155b7f989240cc5cb309e91c68e61536b623e192afd2de138e3f69b3eed730d574bfda9236906004016117da565b60405161158161014a845192602081818801958087835e81015f81520301902054611914565b60405160208185518085835e8101600281520301902060018060a01b0333165f526020526115b560ff60405f205416611b15565b60405160208185518085835e8101600381520301902060018060a01b0333165f526020526115e68260405f20611b77565b6020604051809285518091835e8101600481520301902060018060a01b0333165f5260205260405f206116198154612076565b905551906105ff604051928392339084612084565b3461047f5760208061163f36611799565b9290604051928184925191829101835e810160028152030190209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b3461047f5760208061168c36611799565b9290604051928184925191829101835e810160048152030190209060018060a01b03165f52602052602060405f2054604051908152f35b606081019081106001600160401b038211176106ab57604052565b60c081019081106001600160401b038211176106ab57604052565b90601f801991011681019081106001600160401b038211176106ab57604052565b81601f8201121561047f576020813591016001600160401b0382116106ab5760405192611751601f8401601f1916602001856116f9565b8284528282011161047f57815f92602092838601378301015290565b602435906001600160a01b038216820361047f57565b600435906001600160a01b038216820361047f57565b604060031982011261047f57600435906001600160401b03821161047f576117c39160040161171a565b906024356001600160a01b038116810361047f5790565b919060c08382031261047f57604051906117f3826116de565b819380356001600160401b03811161047f578261181191830161171a565b835260208101356001600160401b03811161047f578261183291830161171a565b602084015260408101356001600160401b03811161047f578261185691830161171a565b60408401526060810135801515810361047f57606084015260808101356001600160401b03811161047f578261188d91830161171a565b608084015260a0810135916001600160401b03831161047f5760a0926118b3920161171a565b910152565b6001600160401b0381116106ab5760051b60200190565b6006548110156118e75760065f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b80548210156118e7575f52600660205f20910201905f90565b90600182811c92168015611942575b602083101461192e57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691611923565b9060405191825f82549261195f84611914565b80845293600181169081156119ca5750600114611986575b50611984925003836116f9565b565b90505f9291925260205f20905f915b8183106119ae575050906020611984928201015f611977565b6020919350806001915483858901015201910190918492611995565b90506020925061198494915060ff191682840152151560051b8201015f611977565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9080602083519182815201916020808360051b8301019401925f915b838310611a3b57505050505090565b9091929394602080611ac4600193601f1986820301875289519060a0611ab3611a95611a83611a73865160c0875260c08701906119ec565b888701518682038a8801526119ec565b604086015185820360408701526119ec565b606085015115156060850152608085015184820360808601526119ec565b9201519060a08184039101526119ec565b97019301930191939290611a2c565b15611ada57565b60405162461bcd60e51b8152602060048201526013602482015272149bdbdb48191bd95cc81b9bdd08195e1a5cdd606a1b6044820152606490fd5b15611b1c57565b60405162461bcd60e51b815260206004820152601860248201527f596f7520617265206e6f7420696e207468697320726f6f6d00000000000000006044820152606490fd5b818110611b6c575050565b5f8155600101611b61565b8054600160401b8110156106ab57611b94916001820181556118fb565b9290926103a35781519283516001600160401b0381116106ab57611bb88254611914565b601f8111612046575b50602094601f8211600114611fe557611bf39293949582915f92611d8d5750508160011b915f199060031b1c19161790565b81555b6001810160208401518051906001600160401b0382116106ab57611c1a8354611914565b601f8111611fb5575b50602090601f8311600114611f5257611c5292915f9183611d8d5750508160011b915f199060031b1c19161790565b90555b6002810160408401518051906001600160401b0382116106ab57611c798354611914565b601f8111611f22575b50602090601f8311600114611ebf57611cb192915f9183611d8d5750508160011b915f199060031b1c19161790565b90555b600381016060840151151560ff801983541691161790556004810160808401518051906001600160401b0382116106ab57611cef8354611914565b601f8111611e8f575b50602090601f8311600114611e26578260a0959360059593611d2e935f92611d8d5750508160011b915f199060031b1c19161790565b90555b019201519182516001600160401b0381116106ab57611d508254611914565b601f8111611df6575b506020601f8211600114611d98578190611d899394955f92611d8d5750508160011b915f199060031b1c19161790565b9055565b015190505f80610b0d565b601f19821690835f52805f20915f5b818110611dde57509583600195969710611dc6575b505050811b019055565b01515f1960f88460031b161c191690555f8080611dbc565b9192602060018192868b015181550194019201611da7565b611e2090835f5260205f20601f840160051c81019160208510610bbd57601f0160051c0190611b61565b5f611d59565b90601f19831691845f52815f20925f5b818110611e77575092600192859260a098966005989610611e5f575b505050811b019055611d31565b01515f1960f88460031b161c191690555f8080611e52565b92936020600181928786015181550195019301611e36565b611eb990845f5260205f20601f850160051c81019160208610610bbd57601f0160051c0190611b61565b5f611cf8565b90601f19831691845f52815f20925f5b818110611f0a5750908460019594939210611ef2575b505050811b019055611cb4565b01515f1960f88460031b161c191690555f8080611ee5565b92936020600181928786015181550195019301611ecf565b611f4c90845f5260205f20601f850160051c81019160208610610bbd57601f0160051c0190611b61565b5f611c82565b90601f19831691845f52815f20925f5b818110611f9d5750908460019594939210611f85575b505050811b019055611c55565b01515f1960f88460031b161c191690555f8080611f78565b92936020600181928786015181550195019301611f62565b611fdf90845f5260205f20601f850160051c81019160208610610bbd57601f0160051c0190611b61565b5f611c23565b601f19821695835f52805f20915f5b88811061202e57508360019596979810612016575b505050811b018155611bf6565b01515f1960f88460031b161c191690555f8080612009565b91926020600181928685015181550194019201611ff4565b61207090835f5260205f20601f840160051c81019160208510610bbd57601f0160051c0190611b61565b5f611bc1565b5f19811461038f5760010190565b61209a6120b894926060835260608301906119ec565b6001600160a01b0390931660208201528083036040909101526119ec565b90565b80548210156118e7575f52600360205f20910201905f90565b909291926103a357825181546001600160a01b0319166001600160a01b03919091161781556020830151805160018301916001600160401b0382116106ab5761211d8354611914565b601f8111612295575b50602090601f831160011461222c5782604095936002959361215c935f92611d8d5750508160011b915f199060031b1c19161790565b90555b019201519182516001600160401b0381116106ab5761217e8254611914565b601f81116121fc575b506020601f82116001146121b7578190611d899394955f92611d8d5750508160011b915f199060031b1c19161790565b601f19821690835f52805f20915f5b8181106121e457509583600195969710611dc657505050811b019055565b9192602060018192868b0151815501940192016121c6565b61222690835f5260205f20601f840160051c81019160208510610bbd57601f0160051c0190611b61565b5f612187565b90601f19831691845f52815f20925f5b81811061227d5750926001928592604098966002989610612265575b505050811b01905561215f565b01515f1960f88460031b161c191690555f8080612258565b9293602060018192878601518155019501930161223c565b6122bf90845f5260205f20601f850160051c81019160208610610bbd57601f0160051c0190611b61565b5f612126565b156122cc57565b60405162461bcd60e51b815260206004820152601760248201527f5061727469636970616e74206e6f7420696e20726f6f6d0000000000000000006044820152606490fd5b9060405161231e816116c3565b60406118b36002839560018060a01b03815416855261233f6001820161194c565b60208601520161194c565b6123548154611914565b908161235e575050565b81601f5f931160011461236f575055565b8183526020832061238b91601f0160051c810190600101611b61565b808252816020812091555556fea2646970667358221220e556ea50cc33c1229cb2fcf0b1d4eecd492d4d23a137bec674f272c6076f7a8064736f6c634300081d0033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedBackends is a free data retrieval call binding the contract method 0x47100107.
//
// Solidity: function authorizedBackends(uint256 ) view returns(address)
func (_Contract *ContractCaller) AuthorizedBackends(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "authorizedBackends", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AuthorizedBackends is a free data retrieval call binding the contract method 0x47100107.
//
// Solidity: function authorizedBackends(uint256 ) view returns(address)
func (_Contract *ContractSession) AuthorizedBackends(arg0 *big.Int) (common.Address, error) {
	return _Contract.Contract.AuthorizedBackends(&_Contract.CallOpts, arg0)
}

// AuthorizedBackends is a free data retrieval call binding the contract method 0x47100107.
//
// Solidity: function authorizedBackends(uint256 ) view returns(address)
func (_Contract *ContractCallerSession) AuthorizedBackends(arg0 *big.Int) (common.Address, error) {
	return _Contract.Contract.AuthorizedBackends(&_Contract.CallOpts, arg0)
}

// GetParticipantInfo is a free data retrieval call binding the contract method 0xe970cf2c.
//
// Solidity: function getParticipantInfo(string _roomId) view returns((address,string,string))
func (_Contract *ContractCaller) GetParticipantInfo(opts *bind.CallOpts, _roomId string) (DAppMeetingParticipant, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getParticipantInfo", _roomId)

	if err != nil {
		return *new(DAppMeetingParticipant), err
	}

	out0 := *abi.ConvertType(out[0], new(DAppMeetingParticipant)).(*DAppMeetingParticipant)

	return out0, err

}

// GetParticipantInfo is a free data retrieval call binding the contract method 0xe970cf2c.
//
// Solidity: function getParticipantInfo(string _roomId) view returns((address,string,string))
func (_Contract *ContractSession) GetParticipantInfo(_roomId string) (DAppMeetingParticipant, error) {
	return _Contract.Contract.GetParticipantInfo(&_Contract.CallOpts, _roomId)
}

// GetParticipantInfo is a free data retrieval call binding the contract method 0xe970cf2c.
//
// Solidity: function getParticipantInfo(string _roomId) view returns((address,string,string))
func (_Contract *ContractCallerSession) GetParticipantInfo(_roomId string) (DAppMeetingParticipant, error) {
	return _Contract.Contract.GetParticipantInfo(&_Contract.CallOpts, _roomId)
}

// GetParticipantTracks is a free data retrieval call binding the contract method 0xc3241b95.
//
// Solidity: function getParticipantTracks(string _roomId, address _participant) view returns((string,string,string,bool,string,string)[])
func (_Contract *ContractCaller) GetParticipantTracks(opts *bind.CallOpts, _roomId string, _participant common.Address) ([]DAppMeetingTrack, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getParticipantTracks", _roomId, _participant)

	if err != nil {
		return *new([]DAppMeetingTrack), err
	}

	out0 := *abi.ConvertType(out[0], new([]DAppMeetingTrack)).(*[]DAppMeetingTrack)

	return out0, err

}

// GetParticipantTracks is a free data retrieval call binding the contract method 0xc3241b95.
//
// Solidity: function getParticipantTracks(string _roomId, address _participant) view returns((string,string,string,bool,string,string)[])
func (_Contract *ContractSession) GetParticipantTracks(_roomId string, _participant common.Address) ([]DAppMeetingTrack, error) {
	return _Contract.Contract.GetParticipantTracks(&_Contract.CallOpts, _roomId, _participant)
}

// GetParticipantTracks is a free data retrieval call binding the contract method 0xc3241b95.
//
// Solidity: function getParticipantTracks(string _roomId, address _participant) view returns((string,string,string,bool,string,string)[])
func (_Contract *ContractCallerSession) GetParticipantTracks(_roomId string, _participant common.Address) ([]DAppMeetingTrack, error) {
	return _Contract.Contract.GetParticipantTracks(&_Contract.CallOpts, _roomId, _participant)
}

// GetParticipantTracksCount is a free data retrieval call binding the contract method 0xf7cc8724.
//
// Solidity: function getParticipantTracksCount(string _roomId, address _participant) view returns(uint256)
func (_Contract *ContractCaller) GetParticipantTracksCount(opts *bind.CallOpts, _roomId string, _participant common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getParticipantTracksCount", _roomId, _participant)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetParticipantTracksCount is a free data retrieval call binding the contract method 0xf7cc8724.
//
// Solidity: function getParticipantTracksCount(string _roomId, address _participant) view returns(uint256)
func (_Contract *ContractSession) GetParticipantTracksCount(_roomId string, _participant common.Address) (*big.Int, error) {
	return _Contract.Contract.GetParticipantTracksCount(&_Contract.CallOpts, _roomId, _participant)
}

// GetParticipantTracksCount is a free data retrieval call binding the contract method 0xf7cc8724.
//
// Solidity: function getParticipantTracksCount(string _roomId, address _participant) view returns(uint256)
func (_Contract *ContractCallerSession) GetParticipantTracksCount(_roomId string, _participant common.Address) (*big.Int, error) {
	return _Contract.Contract.GetParticipantTracksCount(&_Contract.CallOpts, _roomId, _participant)
}

// GetRoomParticipantsCount is a free data retrieval call binding the contract method 0x3f2cc59f.
//
// Solidity: function getRoomParticipantsCount(string _roomId) view returns(uint256)
func (_Contract *ContractCaller) GetRoomParticipantsCount(opts *bind.CallOpts, _roomId string) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getRoomParticipantsCount", _roomId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoomParticipantsCount is a free data retrieval call binding the contract method 0x3f2cc59f.
//
// Solidity: function getRoomParticipantsCount(string _roomId) view returns(uint256)
func (_Contract *ContractSession) GetRoomParticipantsCount(_roomId string) (*big.Int, error) {
	return _Contract.Contract.GetRoomParticipantsCount(&_Contract.CallOpts, _roomId)
}

// GetRoomParticipantsCount is a free data retrieval call binding the contract method 0x3f2cc59f.
//
// Solidity: function getRoomParticipantsCount(string _roomId) view returns(uint256)
func (_Contract *ContractCallerSession) GetRoomParticipantsCount(_roomId string) (*big.Int, error) {
	return _Contract.Contract.GetRoomParticipantsCount(&_Contract.CallOpts, _roomId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCallerSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// ParticipantIndices is a free data retrieval call binding the contract method 0x8df58945.
//
// Solidity: function participantIndices(string , address ) view returns(uint256)
func (_Contract *ContractCaller) ParticipantIndices(opts *bind.CallOpts, arg0 string, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "participantIndices", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParticipantIndices is a free data retrieval call binding the contract method 0x8df58945.
//
// Solidity: function participantIndices(string , address ) view returns(uint256)
func (_Contract *ContractSession) ParticipantIndices(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Contract.Contract.ParticipantIndices(&_Contract.CallOpts, arg0, arg1)
}

// ParticipantIndices is a free data retrieval call binding the contract method 0x8df58945.
//
// Solidity: function participantIndices(string , address ) view returns(uint256)
func (_Contract *ContractCallerSession) ParticipantIndices(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Contract.Contract.ParticipantIndices(&_Contract.CallOpts, arg0, arg1)
}

// ParticipantTrackCount is a free data retrieval call binding the contract method 0x0062b748.
//
// Solidity: function participantTrackCount(string , address ) view returns(uint256)
func (_Contract *ContractCaller) ParticipantTrackCount(opts *bind.CallOpts, arg0 string, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "participantTrackCount", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParticipantTrackCount is a free data retrieval call binding the contract method 0x0062b748.
//
// Solidity: function participantTrackCount(string , address ) view returns(uint256)
func (_Contract *ContractSession) ParticipantTrackCount(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Contract.Contract.ParticipantTrackCount(&_Contract.CallOpts, arg0, arg1)
}

// ParticipantTrackCount is a free data retrieval call binding the contract method 0x0062b748.
//
// Solidity: function participantTrackCount(string , address ) view returns(uint256)
func (_Contract *ContractCallerSession) ParticipantTrackCount(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Contract.Contract.ParticipantTrackCount(&_Contract.CallOpts, arg0, arg1)
}

// ParticipantTracks is a free data retrieval call binding the contract method 0x5b85facc.
//
// Solidity: function participantTracks(string , address , uint256 ) view returns(string trackName, string mid, string location, bool isPublished, string sessionId, string roomId)
func (_Contract *ContractCaller) ParticipantTracks(opts *bind.CallOpts, arg0 string, arg1 common.Address, arg2 *big.Int) (struct {
	TrackName   string
	Mid         string
	Location    string
	IsPublished bool
	SessionId   string
	RoomId      string
}, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "participantTracks", arg0, arg1, arg2)

	outstruct := new(struct {
		TrackName   string
		Mid         string
		Location    string
		IsPublished bool
		SessionId   string
		RoomId      string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TrackName = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Mid = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Location = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.IsPublished = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.SessionId = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.RoomId = *abi.ConvertType(out[5], new(string)).(*string)

	return *outstruct, err

}

// ParticipantTracks is a free data retrieval call binding the contract method 0x5b85facc.
//
// Solidity: function participantTracks(string , address , uint256 ) view returns(string trackName, string mid, string location, bool isPublished, string sessionId, string roomId)
func (_Contract *ContractSession) ParticipantTracks(arg0 string, arg1 common.Address, arg2 *big.Int) (struct {
	TrackName   string
	Mid         string
	Location    string
	IsPublished bool
	SessionId   string
	RoomId      string
}, error) {
	return _Contract.Contract.ParticipantTracks(&_Contract.CallOpts, arg0, arg1, arg2)
}

// ParticipantTracks is a free data retrieval call binding the contract method 0x5b85facc.
//
// Solidity: function participantTracks(string , address , uint256 ) view returns(string trackName, string mid, string location, bool isPublished, string sessionId, string roomId)
func (_Contract *ContractCallerSession) ParticipantTracks(arg0 string, arg1 common.Address, arg2 *big.Int) (struct {
	TrackName   string
	Mid         string
	Location    string
	IsPublished bool
	SessionId   string
	RoomId      string
}, error) {
	return _Contract.Contract.ParticipantTracks(&_Contract.CallOpts, arg0, arg1, arg2)
}

// ParticipantsInRoom is a free data retrieval call binding the contract method 0x2aba9dda.
//
// Solidity: function participantsInRoom(string , address ) view returns(bool)
func (_Contract *ContractCaller) ParticipantsInRoom(opts *bind.CallOpts, arg0 string, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "participantsInRoom", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ParticipantsInRoom is a free data retrieval call binding the contract method 0x2aba9dda.
//
// Solidity: function participantsInRoom(string , address ) view returns(bool)
func (_Contract *ContractSession) ParticipantsInRoom(arg0 string, arg1 common.Address) (bool, error) {
	return _Contract.Contract.ParticipantsInRoom(&_Contract.CallOpts, arg0, arg1)
}

// ParticipantsInRoom is a free data retrieval call binding the contract method 0x2aba9dda.
//
// Solidity: function participantsInRoom(string , address ) view returns(bool)
func (_Contract *ContractCallerSession) ParticipantsInRoom(arg0 string, arg1 common.Address) (bool, error) {
	return _Contract.Contract.ParticipantsInRoom(&_Contract.CallOpts, arg0, arg1)
}

// Rooms is a free data retrieval call binding the contract method 0xbbbcc869.
//
// Solidity: function rooms(string ) view returns(string roomId, uint256 creationTime)
func (_Contract *ContractCaller) Rooms(opts *bind.CallOpts, arg0 string) (struct {
	RoomId       string
	CreationTime *big.Int
}, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "rooms", arg0)

	outstruct := new(struct {
		RoomId       string
		CreationTime *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoomId = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.CreationTime = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Rooms is a free data retrieval call binding the contract method 0xbbbcc869.
//
// Solidity: function rooms(string ) view returns(string roomId, uint256 creationTime)
func (_Contract *ContractSession) Rooms(arg0 string) (struct {
	RoomId       string
	CreationTime *big.Int
}, error) {
	return _Contract.Contract.Rooms(&_Contract.CallOpts, arg0)
}

// Rooms is a free data retrieval call binding the contract method 0xbbbcc869.
//
// Solidity: function rooms(string ) view returns(string roomId, uint256 creationTime)
func (_Contract *ContractCallerSession) Rooms(arg0 string) (struct {
	RoomId       string
	CreationTime *big.Int
}, error) {
	return _Contract.Contract.Rooms(&_Contract.CallOpts, arg0)
}

// AddAuthorizedBackend is a paid mutator transaction binding the contract method 0xf2776d2e.
//
// Solidity: function addAuthorizedBackend(address _backend) returns()
func (_Contract *ContractTransactor) AddAuthorizedBackend(opts *bind.TransactOpts, _backend common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addAuthorizedBackend", _backend)
}

// AddAuthorizedBackend is a paid mutator transaction binding the contract method 0xf2776d2e.
//
// Solidity: function addAuthorizedBackend(address _backend) returns()
func (_Contract *ContractSession) AddAuthorizedBackend(_backend common.Address) (*types.Transaction, error) {
	return _Contract.Contract.AddAuthorizedBackend(&_Contract.TransactOpts, _backend)
}

// AddAuthorizedBackend is a paid mutator transaction binding the contract method 0xf2776d2e.
//
// Solidity: function addAuthorizedBackend(address _backend) returns()
func (_Contract *ContractTransactorSession) AddAuthorizedBackend(_backend common.Address) (*types.Transaction, error) {
	return _Contract.Contract.AddAuthorizedBackend(&_Contract.TransactOpts, _backend)
}

// AddTrack is a paid mutator transaction binding the contract method 0x3b4ddda8.
//
// Solidity: function addTrack(string _roomId, (string,string,string,bool,string,string) _newTrack) returns()
func (_Contract *ContractTransactor) AddTrack(opts *bind.TransactOpts, _roomId string, _newTrack DAppMeetingTrack) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addTrack", _roomId, _newTrack)
}

// AddTrack is a paid mutator transaction binding the contract method 0x3b4ddda8.
//
// Solidity: function addTrack(string _roomId, (string,string,string,bool,string,string) _newTrack) returns()
func (_Contract *ContractSession) AddTrack(_roomId string, _newTrack DAppMeetingTrack) (*types.Transaction, error) {
	return _Contract.Contract.AddTrack(&_Contract.TransactOpts, _roomId, _newTrack)
}

// AddTrack is a paid mutator transaction binding the contract method 0x3b4ddda8.
//
// Solidity: function addTrack(string _roomId, (string,string,string,bool,string,string) _newTrack) returns()
func (_Contract *ContractTransactorSession) AddTrack(_roomId string, _newTrack DAppMeetingTrack) (*types.Transaction, error) {
	return _Contract.Contract.AddTrack(&_Contract.TransactOpts, _roomId, _newTrack)
}

// CreateRoom is a paid mutator transaction binding the contract method 0x7306d2dd.
//
// Solidity: function createRoom(string _roomId) returns()
func (_Contract *ContractTransactor) CreateRoom(opts *bind.TransactOpts, _roomId string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "createRoom", _roomId)
}

// CreateRoom is a paid mutator transaction binding the contract method 0x7306d2dd.
//
// Solidity: function createRoom(string _roomId) returns()
func (_Contract *ContractSession) CreateRoom(_roomId string) (*types.Transaction, error) {
	return _Contract.Contract.CreateRoom(&_Contract.TransactOpts, _roomId)
}

// CreateRoom is a paid mutator transaction binding the contract method 0x7306d2dd.
//
// Solidity: function createRoom(string _roomId) returns()
func (_Contract *ContractTransactorSession) CreateRoom(_roomId string) (*types.Transaction, error) {
	return _Contract.Contract.CreateRoom(&_Contract.TransactOpts, _roomId)
}

// ForwardEventToBackend is a paid mutator transaction binding the contract method 0xaaea8180.
//
// Solidity: function forwardEventToBackend(string _roomId, bytes _eventData) returns()
func (_Contract *ContractTransactor) ForwardEventToBackend(opts *bind.TransactOpts, _roomId string, _eventData []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "forwardEventToBackend", _roomId, _eventData)
}

// ForwardEventToBackend is a paid mutator transaction binding the contract method 0xaaea8180.
//
// Solidity: function forwardEventToBackend(string _roomId, bytes _eventData) returns()
func (_Contract *ContractSession) ForwardEventToBackend(_roomId string, _eventData []byte) (*types.Transaction, error) {
	return _Contract.Contract.ForwardEventToBackend(&_Contract.TransactOpts, _roomId, _eventData)
}

// ForwardEventToBackend is a paid mutator transaction binding the contract method 0xaaea8180.
//
// Solidity: function forwardEventToBackend(string _roomId, bytes _eventData) returns()
func (_Contract *ContractTransactorSession) ForwardEventToBackend(_roomId string, _eventData []byte) (*types.Transaction, error) {
	return _Contract.Contract.ForwardEventToBackend(&_Contract.TransactOpts, _roomId, _eventData)
}

// ForwardEventToFrontend is a paid mutator transaction binding the contract method 0xf3398cf6.
//
// Solidity: function forwardEventToFrontend(string _roomId, address _participant, bytes _eventData) returns()
func (_Contract *ContractTransactor) ForwardEventToFrontend(opts *bind.TransactOpts, _roomId string, _participant common.Address, _eventData []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "forwardEventToFrontend", _roomId, _participant, _eventData)
}

// ForwardEventToFrontend is a paid mutator transaction binding the contract method 0xf3398cf6.
//
// Solidity: function forwardEventToFrontend(string _roomId, address _participant, bytes _eventData) returns()
func (_Contract *ContractSession) ForwardEventToFrontend(_roomId string, _participant common.Address, _eventData []byte) (*types.Transaction, error) {
	return _Contract.Contract.ForwardEventToFrontend(&_Contract.TransactOpts, _roomId, _participant, _eventData)
}

// ForwardEventToFrontend is a paid mutator transaction binding the contract method 0xf3398cf6.
//
// Solidity: function forwardEventToFrontend(string _roomId, address _participant, bytes _eventData) returns()
func (_Contract *ContractTransactorSession) ForwardEventToFrontend(_roomId string, _participant common.Address, _eventData []byte) (*types.Transaction, error) {
	return _Contract.Contract.ForwardEventToFrontend(&_Contract.TransactOpts, _roomId, _participant, _eventData)
}

// JoinRoom is a paid mutator transaction binding the contract method 0x3c229779.
//
// Solidity: function joinRoom(string _roomId, string _name, (string,string,string,bool,string,string)[] _initialTracks, bytes sessionDescription) returns()
func (_Contract *ContractTransactor) JoinRoom(opts *bind.TransactOpts, _roomId string, _name string, _initialTracks []DAppMeetingTrack, sessionDescription []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "joinRoom", _roomId, _name, _initialTracks, sessionDescription)
}

// JoinRoom is a paid mutator transaction binding the contract method 0x3c229779.
//
// Solidity: function joinRoom(string _roomId, string _name, (string,string,string,bool,string,string)[] _initialTracks, bytes sessionDescription) returns()
func (_Contract *ContractSession) JoinRoom(_roomId string, _name string, _initialTracks []DAppMeetingTrack, sessionDescription []byte) (*types.Transaction, error) {
	return _Contract.Contract.JoinRoom(&_Contract.TransactOpts, _roomId, _name, _initialTracks, sessionDescription)
}

// JoinRoom is a paid mutator transaction binding the contract method 0x3c229779.
//
// Solidity: function joinRoom(string _roomId, string _name, (string,string,string,bool,string,string)[] _initialTracks, bytes sessionDescription) returns()
func (_Contract *ContractTransactorSession) JoinRoom(_roomId string, _name string, _initialTracks []DAppMeetingTrack, sessionDescription []byte) (*types.Transaction, error) {
	return _Contract.Contract.JoinRoom(&_Contract.TransactOpts, _roomId, _name, _initialTracks, sessionDescription)
}

// LeaveRoom is a paid mutator transaction binding the contract method 0xfab38543.
//
// Solidity: function leaveRoom(string _roomId) returns()
func (_Contract *ContractTransactor) LeaveRoom(opts *bind.TransactOpts, _roomId string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "leaveRoom", _roomId)
}

// LeaveRoom is a paid mutator transaction binding the contract method 0xfab38543.
//
// Solidity: function leaveRoom(string _roomId) returns()
func (_Contract *ContractSession) LeaveRoom(_roomId string) (*types.Transaction, error) {
	return _Contract.Contract.LeaveRoom(&_Contract.TransactOpts, _roomId)
}

// LeaveRoom is a paid mutator transaction binding the contract method 0xfab38543.
//
// Solidity: function leaveRoom(string _roomId) returns()
func (_Contract *ContractTransactorSession) LeaveRoom(_roomId string) (*types.Transaction, error) {
	return _Contract.Contract.LeaveRoom(&_Contract.TransactOpts, _roomId)
}

// RemoveAuthorizedBackend is a paid mutator transaction binding the contract method 0x9ba78fd9.
//
// Solidity: function removeAuthorizedBackend(address _backend) returns()
func (_Contract *ContractTransactor) RemoveAuthorizedBackend(opts *bind.TransactOpts, _backend common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "removeAuthorizedBackend", _backend)
}

// RemoveAuthorizedBackend is a paid mutator transaction binding the contract method 0x9ba78fd9.
//
// Solidity: function removeAuthorizedBackend(address _backend) returns()
func (_Contract *ContractSession) RemoveAuthorizedBackend(_backend common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RemoveAuthorizedBackend(&_Contract.TransactOpts, _backend)
}

// RemoveAuthorizedBackend is a paid mutator transaction binding the contract method 0x9ba78fd9.
//
// Solidity: function removeAuthorizedBackend(address _backend) returns()
func (_Contract *ContractTransactorSession) RemoveAuthorizedBackend(_backend common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RemoveAuthorizedBackend(&_Contract.TransactOpts, _backend)
}

// SetParticipantSessionID is a paid mutator transaction binding the contract method 0xbb9af00d.
//
// Solidity: function setParticipantSessionID(string _roomId, address _participantAddress, string _sessionID) returns()
func (_Contract *ContractTransactor) SetParticipantSessionID(opts *bind.TransactOpts, _roomId string, _participantAddress common.Address, _sessionID string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setParticipantSessionID", _roomId, _participantAddress, _sessionID)
}

// SetParticipantSessionID is a paid mutator transaction binding the contract method 0xbb9af00d.
//
// Solidity: function setParticipantSessionID(string _roomId, address _participantAddress, string _sessionID) returns()
func (_Contract *ContractSession) SetParticipantSessionID(_roomId string, _participantAddress common.Address, _sessionID string) (*types.Transaction, error) {
	return _Contract.Contract.SetParticipantSessionID(&_Contract.TransactOpts, _roomId, _participantAddress, _sessionID)
}

// SetParticipantSessionID is a paid mutator transaction binding the contract method 0xbb9af00d.
//
// Solidity: function setParticipantSessionID(string _roomId, address _participantAddress, string _sessionID) returns()
func (_Contract *ContractTransactorSession) SetParticipantSessionID(_roomId string, _participantAddress common.Address, _sessionID string) (*types.Transaction, error) {
	return _Contract.Contract.SetParticipantSessionID(&_Contract.TransactOpts, _roomId, _participantAddress, _sessionID)
}

// ContractEventForwardedToBackendIterator is returned from FilterEventForwardedToBackend and is used to iterate over the raw logs and unpacked data for EventForwardedToBackend events raised by the Contract contract.
type ContractEventForwardedToBackendIterator struct {
	Event *ContractEventForwardedToBackend // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractEventForwardedToBackendIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEventForwardedToBackend)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractEventForwardedToBackend)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractEventForwardedToBackendIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEventForwardedToBackendIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEventForwardedToBackend represents a EventForwardedToBackend event raised by the Contract contract.
type ContractEventForwardedToBackend struct {
	RoomId    string
	Sender    common.Address
	EventData []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEventForwardedToBackend is a free log retrieval operation binding the contract event 0x5af2a6a9c8117113b24e80f03488d6bb6d93973b5faed013eca8657949140f1e.
//
// Solidity: event EventForwardedToBackend(string roomId, address sender, bytes eventData)
func (_Contract *ContractFilterer) FilterEventForwardedToBackend(opts *bind.FilterOpts) (*ContractEventForwardedToBackendIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "EventForwardedToBackend")
	if err != nil {
		return nil, err
	}
	return &ContractEventForwardedToBackendIterator{contract: _Contract.contract, event: "EventForwardedToBackend", logs: logs, sub: sub}, nil
}

// WatchEventForwardedToBackend is a free log subscription operation binding the contract event 0x5af2a6a9c8117113b24e80f03488d6bb6d93973b5faed013eca8657949140f1e.
//
// Solidity: event EventForwardedToBackend(string roomId, address sender, bytes eventData)
func (_Contract *ContractFilterer) WatchEventForwardedToBackend(opts *bind.WatchOpts, sink chan<- *ContractEventForwardedToBackend) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "EventForwardedToBackend")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEventForwardedToBackend)
				if err := _Contract.contract.UnpackLog(event, "EventForwardedToBackend", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEventForwardedToBackend is a log parse operation binding the contract event 0x5af2a6a9c8117113b24e80f03488d6bb6d93973b5faed013eca8657949140f1e.
//
// Solidity: event EventForwardedToBackend(string roomId, address sender, bytes eventData)
func (_Contract *ContractFilterer) ParseEventForwardedToBackend(log types.Log) (*ContractEventForwardedToBackend, error) {
	event := new(ContractEventForwardedToBackend)
	if err := _Contract.contract.UnpackLog(event, "EventForwardedToBackend", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEventForwardedToFrontendIterator is returned from FilterEventForwardedToFrontend and is used to iterate over the raw logs and unpacked data for EventForwardedToFrontend events raised by the Contract contract.
type ContractEventForwardedToFrontendIterator struct {
	Event *ContractEventForwardedToFrontend // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractEventForwardedToFrontendIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEventForwardedToFrontend)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractEventForwardedToFrontend)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractEventForwardedToFrontendIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEventForwardedToFrontendIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEventForwardedToFrontend represents a EventForwardedToFrontend event raised by the Contract contract.
type ContractEventForwardedToFrontend struct {
	RoomId      string
	Participant common.Address
	EventData   []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEventForwardedToFrontend is a free log retrieval operation binding the contract event 0xb588e6e55cc56c9ac78812d33d4f8fbbebe720ee569bed8fd4615ddb178fcb12.
//
// Solidity: event EventForwardedToFrontend(string roomId, address participant, bytes eventData)
func (_Contract *ContractFilterer) FilterEventForwardedToFrontend(opts *bind.FilterOpts) (*ContractEventForwardedToFrontendIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "EventForwardedToFrontend")
	if err != nil {
		return nil, err
	}
	return &ContractEventForwardedToFrontendIterator{contract: _Contract.contract, event: "EventForwardedToFrontend", logs: logs, sub: sub}, nil
}

// WatchEventForwardedToFrontend is a free log subscription operation binding the contract event 0xb588e6e55cc56c9ac78812d33d4f8fbbebe720ee569bed8fd4615ddb178fcb12.
//
// Solidity: event EventForwardedToFrontend(string roomId, address participant, bytes eventData)
func (_Contract *ContractFilterer) WatchEventForwardedToFrontend(opts *bind.WatchOpts, sink chan<- *ContractEventForwardedToFrontend) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "EventForwardedToFrontend")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEventForwardedToFrontend)
				if err := _Contract.contract.UnpackLog(event, "EventForwardedToFrontend", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEventForwardedToFrontend is a log parse operation binding the contract event 0xb588e6e55cc56c9ac78812d33d4f8fbbebe720ee569bed8fd4615ddb178fcb12.
//
// Solidity: event EventForwardedToFrontend(string roomId, address participant, bytes eventData)
func (_Contract *ContractFilterer) ParseEventForwardedToFrontend(log types.Log) (*ContractEventForwardedToFrontend, error) {
	event := new(ContractEventForwardedToFrontend)
	if err := _Contract.contract.UnpackLog(event, "EventForwardedToFrontend", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractParticipantJoinedIterator is returned from FilterParticipantJoined and is used to iterate over the raw logs and unpacked data for ParticipantJoined events raised by the Contract contract.
type ContractParticipantJoinedIterator struct {
	Event *ContractParticipantJoined // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractParticipantJoinedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractParticipantJoined)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractParticipantJoined)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractParticipantJoinedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractParticipantJoinedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractParticipantJoined represents a ParticipantJoined event raised by the Contract contract.
type ContractParticipantJoined struct {
	RoomId             string
	Participant        common.Address
	InitialTracks      []DAppMeetingTrack
	SessionDescription []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterParticipantJoined is a free log retrieval operation binding the contract event 0x21e27f169d16f9bc6b07bf8a2c343f07f22bfd14122a3181cf25419574f0109b.
//
// Solidity: event ParticipantJoined(string roomId, address participant, (string,string,string,bool,string,string)[] initialTracks, bytes sessionDescription)
func (_Contract *ContractFilterer) FilterParticipantJoined(opts *bind.FilterOpts) (*ContractParticipantJoinedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ParticipantJoined")
	if err != nil {
		return nil, err
	}
	return &ContractParticipantJoinedIterator{contract: _Contract.contract, event: "ParticipantJoined", logs: logs, sub: sub}, nil
}

// WatchParticipantJoined is a free log subscription operation binding the contract event 0x21e27f169d16f9bc6b07bf8a2c343f07f22bfd14122a3181cf25419574f0109b.
//
// Solidity: event ParticipantJoined(string roomId, address participant, (string,string,string,bool,string,string)[] initialTracks, bytes sessionDescription)
func (_Contract *ContractFilterer) WatchParticipantJoined(opts *bind.WatchOpts, sink chan<- *ContractParticipantJoined) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ParticipantJoined")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractParticipantJoined)
				if err := _Contract.contract.UnpackLog(event, "ParticipantJoined", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseParticipantJoined is a log parse operation binding the contract event 0x21e27f169d16f9bc6b07bf8a2c343f07f22bfd14122a3181cf25419574f0109b.
//
// Solidity: event ParticipantJoined(string roomId, address participant, (string,string,string,bool,string,string)[] initialTracks, bytes sessionDescription)
func (_Contract *ContractFilterer) ParseParticipantJoined(log types.Log) (*ContractParticipantJoined, error) {
	event := new(ContractParticipantJoined)
	if err := _Contract.contract.UnpackLog(event, "ParticipantJoined", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractParticipantLeftIterator is returned from FilterParticipantLeft and is used to iterate over the raw logs and unpacked data for ParticipantLeft events raised by the Contract contract.
type ContractParticipantLeftIterator struct {
	Event *ContractParticipantLeft // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractParticipantLeftIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractParticipantLeft)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractParticipantLeft)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractParticipantLeftIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractParticipantLeftIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractParticipantLeft represents a ParticipantLeft event raised by the Contract contract.
type ContractParticipantLeft struct {
	RoomId      string
	Participant common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterParticipantLeft is a free log retrieval operation binding the contract event 0xb4d5a3866c2b36a076b8cac89e6068deeff68398781315f07060bdead654e63b.
//
// Solidity: event ParticipantLeft(string roomId, address participant)
func (_Contract *ContractFilterer) FilterParticipantLeft(opts *bind.FilterOpts) (*ContractParticipantLeftIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ParticipantLeft")
	if err != nil {
		return nil, err
	}
	return &ContractParticipantLeftIterator{contract: _Contract.contract, event: "ParticipantLeft", logs: logs, sub: sub}, nil
}

// WatchParticipantLeft is a free log subscription operation binding the contract event 0xb4d5a3866c2b36a076b8cac89e6068deeff68398781315f07060bdead654e63b.
//
// Solidity: event ParticipantLeft(string roomId, address participant)
func (_Contract *ContractFilterer) WatchParticipantLeft(opts *bind.WatchOpts, sink chan<- *ContractParticipantLeft) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ParticipantLeft")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractParticipantLeft)
				if err := _Contract.contract.UnpackLog(event, "ParticipantLeft", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseParticipantLeft is a log parse operation binding the contract event 0xb4d5a3866c2b36a076b8cac89e6068deeff68398781315f07060bdead654e63b.
//
// Solidity: event ParticipantLeft(string roomId, address participant)
func (_Contract *ContractFilterer) ParseParticipantLeft(log types.Log) (*ContractParticipantLeft, error) {
	event := new(ContractParticipantLeft)
	if err := _Contract.contract.UnpackLog(event, "ParticipantLeft", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractTrackAddedIterator is returned from FilterTrackAdded and is used to iterate over the raw logs and unpacked data for TrackAdded events raised by the Contract contract.
type ContractTrackAddedIterator struct {
	Event *ContractTrackAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractTrackAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractTrackAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractTrackAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractTrackAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractTrackAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractTrackAdded represents a TrackAdded event raised by the Contract contract.
type ContractTrackAdded struct {
	RoomId      string
	Participant common.Address
	TrackName   string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTrackAdded is a free log retrieval operation binding the contract event 0x989240cc5cb309e91c68e61536b623e192afd2de138e3f69b3eed730d574bfda.
//
// Solidity: event TrackAdded(string roomId, address participant, string trackName)
func (_Contract *ContractFilterer) FilterTrackAdded(opts *bind.FilterOpts) (*ContractTrackAddedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "TrackAdded")
	if err != nil {
		return nil, err
	}
	return &ContractTrackAddedIterator{contract: _Contract.contract, event: "TrackAdded", logs: logs, sub: sub}, nil
}

// WatchTrackAdded is a free log subscription operation binding the contract event 0x989240cc5cb309e91c68e61536b623e192afd2de138e3f69b3eed730d574bfda.
//
// Solidity: event TrackAdded(string roomId, address participant, string trackName)
func (_Contract *ContractFilterer) WatchTrackAdded(opts *bind.WatchOpts, sink chan<- *ContractTrackAdded) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "TrackAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractTrackAdded)
				if err := _Contract.contract.UnpackLog(event, "TrackAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTrackAdded is a log parse operation binding the contract event 0x989240cc5cb309e91c68e61536b623e192afd2de138e3f69b3eed730d574bfda.
//
// Solidity: event TrackAdded(string roomId, address participant, string trackName)
func (_Contract *ContractFilterer) ParseTrackAdded(log types.Log) (*ContractTrackAdded, error) {
	event := new(ContractTrackAdded)
	if err := _Contract.contract.UnpackLog(event, "TrackAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
