import { Client, Server } from "bittorrent-tracker";
import { readFile, writeFileSync } from "fs";
import WebTorrent from "webtorrent";

const server = new Server({
  udp: false, // enable udp server? [default=true]
  http: true, // enable http server? [default=true]
  ws: false, // enable websocket server? [default=true]
  stats: true, // enable web-based statistics? [default=true]
  trustProxy: false, // enable trusting x-forwarded-for header for remote IP [default=false]
  filter: function (infoHash, params, cb) {
    const allowed = infoHash === "49ed8b48c132974c5a30fc5f7b6e1fe77737a4df";

    if (allowed) {
      // If the callback is passed `null`, the torrent will be allowed.
      cb(null);
    } else {
      // If the callback is passed an `Error` object, the torrent will be disallowed
      // and the error's `message` property will be given as the reason.
      cb(new Error("disallowed torrent"));
    }
  },
});

// Internal http, udp, and websocket servers exposed as public properties.
server.on("error", function (err) {
  // fatal server error!
  console.log(err.message);
});

server.on("warning", function (err) {
  // client sent bad data. probably not a problem, just a buggy client.
  console.log(err.message);
});

server.on("listening", function () {
  // fired when all requested servers are listening

  // HTTP
  const httpAddr = server.http.address();
  const httpHost = httpAddr.address !== "::" ? httpAddr.address : "localhost";
  const httpPort = httpAddr.port;
  console.log(`HTTP tracker: http://${httpHost}:${httpPort}/announce`);
});

// start tracker server listening! Use 0 to listen on a random free port.
const port = 8080;
const hostname = "0.0.0.0";
server.listen(port, hostname, () => {
  // Do something on listening...
});

// listen for individual tracker messages from peers:

server.on("start", function (addr) {
  console.log("got start message from " + addr);
});

server.on("complete", function (addr) {
  console.log("got complete message from " + addr);
});

server.on("update", function (addr) {
  console.log("got update message from " + addr);
});

server.on("stop", function (addr) {
  console.log("got stop message from " + addr);
});

// get info hashes for all torrents in the tracker server
console.log(Object.keys(server.torrents));

const seeder = new Client({
  infoHash: "49ed8b48c132974c5a30fc5f7b6e1fe77737a4df", // hex string or Buffer
  peerId: "aaa67059ed6bd08362da625b3ae77f6f4a075aaa", // hex string or Buffer
  announce: ["http://localhost:8080/announce"], // list of tracker server urls
  port: 6881, // torrent client port, (in browser, optional)
});

const initialSeeder = new WebTorrent({
  // maxConns: Number,        // Max number of connections per torrent (default=55)
  // nodeId: String|Uint8Array,   // DHT protocol node ID (default=randomly generated)
  // peerId: String|Uint8Array,   // Wire protocol peer ID (default=randomly generated)
  // tracker: Boolean|Object, // Enable trackers (default=true), or options object for Tracker
  dht: false, // Enable DHT (default=true), or options object for DHT
  lsd: false, // Enable BEP14 local service discovery (default=true)
  utPex: false, // Enable BEP11 Peer Exchange (default=true)
  // natUpnp:  String, // Enable NAT port mapping via NAT-UPnP (default=true). NodeJS only
  // natPmp: Boolean,         // Enable NAT port mapping via NAT-PMP (default=true). NodeJS only.
  webSeeds: false, // Enable BEP19 web seeds (default=true)
  utp: false, // Enable BEP29 uTorrent transport protocol (default=true)
  // blocklist: Array|String, // List of IP's to block
  // downloadLimit: Number,   // Max download speed (bytes/sec) over all torrents (default=-1)
  // uploadLimit: Number,     // Max upload speed (bytes/sec) over all torrents (default=-1)
});

initialSeeder.seed(
  "./congratulations.gif",
  { announce: ["http://0.0.0.0:8080/announce"] },
  (torrent) => {
    console.log("initiated seeding", torrent.infoHash);
    writeFileSync("./congratulations.gif.torrent", torrent.torrentFile);
    seeder.start();
  }
);
