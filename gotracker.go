package gotracker
import "http"
import "./bencode"
import "fmt"
import "os"
import "container/list"
import "strconv"
import "bufio"

type Tracker struct {
	infoHash string;
	peerId string;
	url string;
	port string;
	interval int;
}

func (t *Tracker) Init(url string, infoHash string,peerId string,port string){
	t.url = url;
	t.infoHash = infoHash;
	t.peerId = peerId;
	t.port = port;
}

func (t *Tracker) Request(uploaded int, downloaded int, left int, status string) (peers *list.List,err os.Error){
	err = nil;
	url:= fmt.Sprint(t.url,
	"?",
	"info_hash=",http.URLEscape(t.infoHash),
	"&peer_id=",http.URLEscape(t.peerId),
	"&port=",http.URLEscape(t.port),
	"&uploaded=",strconv.Itoa(uploaded),
	"&downloaded=",strconv.Itoa(downloaded),
	"&left=",strconv.Itoa(left),
	"&status=",http.URLEscape(status));
	println(url);
	response,_, err := http.Get(url);
	if err != nil { return;}

	if response.StatusCode != http.StatusOK {
		err = os.NewError("http error");
	}
	buf := new(bencode.BeString);
	resReader := bufio.NewReader(response.Body);
	be,err := buf.Decode(resReader);
	if err != nil {
		return;
	}
		if be.Betype != bencode.Bedict {
			err = os.NewError("unexpected response from tracker");
			return;
		}
		if failure,ok:= be.Bedict["failure reason"]; !ok {
			err = os.NewError("unexpected response from tracker");
			return;
		} else {
			print(failure);
			return;
		}
		if interval,ok := be.Bedict["interval"]; !ok {
			err = os.NewError("unexpected response from tracker");
			return;
		} else {
			t.interval,err = strconv.Atoi(interval.Bestr);
		}
		if _,ok := be.Bedict["peers"]; !ok {
			err = os.NewError("unexpected response from tracker");
			return;
		}
		list:=list.New();
		for i := range be.Bedict["peers"].Belist.Iter(){
			if i.(*bencode.BeNode).Betype == bencode.Bedict {
				list.PushFront(i.(*bencode.BeNode).Bedict);
			} else {
				err = os.NewError("unexpected response from tracker");
				return;
			}
		}
		peers = list;
		return;
}
/*
request:

info_hash: urlencoded 20-byte SHA1 hash of the value of the info key from the Metainfo file. Note that the value will be a bencoded dictionary, given the definition of the info key above.

peer_id: urlencoded 20-byte string used as a unique ID for the client, generated by the client at startup. This is allowed to be any value, and may be binary data. There are currently no guidelines for generating this peer ID. However, one may rightly presume that it must at least be unique for your local machine, thus should probably incorporate things like process ID and perhaps a timestamp recorded at startup. See peer_id below for common client encodings of this field.

port: The port number that the client is listening on. Ports reserved for BitTorrent are typically 6881-6889. Clients may choose to give up if it cannot establish a port within this range.

uploaded: The total amount uploaded (since the client sent the 'started' event to the tracker) in base ten ASCII. While not explicitly stated in the official specification, the concensus is that this should be the total number of bytes uploaded.

downloaded: The total amount downloaded (since the client sent the 'started' event to the tracker) in base ten ASCII. While not explicitly stated in the official specification, the consensus is that this should be the total number of bytes downloaded.

left: The number of bytes this client still has to download, encoded in base ten ASCII.

compact: Setting this to 1 indicates that the client accepts a compact response. The peers list is replaced by a peers string with 6 bytes per peer. The first four bytes are the host (in network byte order), the last two bytes are the port (again in network byte order). It should be noted that some trackers only support compact responses (for saving bandwidth) and either refuse requests without "compact=1" or simply send a compact response unless the request contains "compact=0" (in which case they will refuse the request.)

no_peer_id: Indicates that the tracker can omit peer id field in peers dictionary. This option is ignored if compact is enabled.
event: If specified, must be one of started, completed, stopped, (or empty which is the same as not being specified). If not specified, then this request is one performed at regular intervals.

started: The first request to the tracker must include the event key with this value.

stopped: Must be sent to the tracker if the client is shutting down gracefully.

completed: Must be sent to the tracker when the download completes. However, must not be sent if the download was already 100% complete when the client started. Presumably, this is to allow the tracker to increment the "completed downloads" metric based solely on this event.

ip: Optional. The true IP address of the client machine, in dotted quad format or rfc3513 defined hexed IPv6 address. Notes: In general this parameter is not necessary as the address of the client can be determined from the IP address from which the HTTP request came. The parameter is only needed in the case where the IP address that the request came in on is not the IP address of the client. This happens if the client is communicating to the tracker through a proxy (or a transparent web proxy/cache.) It also is necessary when both the client and the tracker are on the same local side of a NAT gateway. The reason for this is that otherwise the tracker would give out the internal (RFC1918) address of the client, which is not routable. Therefore the client must explicitly state its (external, routable) IP address to be given out to external peers. Various trackers treat this parameter differently. Some only honor it only if the IP address that the request came in on is in RFC1918 space. Others honor it unconditionally, while others ignore it completely. In case of IPv6 address (e.g.: 2001:db8:1:2::100) it indicates only that client can communicate via IPv6.

numwant: Optional. Number of peers that the client would like to receive from the tracker. This value is permitted to be zero. If omitted, typically defaults to 50 peers.

key: Optional. An additional identification that is not shared with any users. It is intended to allow a client to prove their identity should their IP address change.

trackerid: Optional. If a previous announce contained a tracker id, it should be set here.



response:

Tracker responses are bencoded dictionaries. If a tracker response has a key failure reason, then that maps to a human readable string which explains why the query failed, and no other keys are required. Otherwise, it must have two keys: interval, which maps to the number of seconds the downloader should wait between regular rerequests, and peers. peers maps to a list of dictionaries corresponding to peers, each of which contains the keys peer id, ip, and port, which map to the peer's self-selected ID, IP address or dns name as a string, and port number, respectively. Note that downloaders may rerequest on nonscheduled times if an event happens or they need more peers.

*/
































































