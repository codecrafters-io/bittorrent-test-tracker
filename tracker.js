import { Server } from "bittorrent-tracker";

const server = new Server({
  udp: false, // enable udp server? [default=true]
  http: true, // enable http server? [default=true]
  ws: false, // enable websocket server? [default=true]
  stats: true, // enable web-based statistics? [default=true]
  trustProxy: false, // enable trusting x-forwarded-for header for remote IP [default=false]
  filter: function (infoHash, params, cb) {
    console.log("filter called", params, infoHash);
    const allowed = infoHash === "49ed8b48c132974c5a30fc5f7b6e1fe77737a4df"; // congratulations.gif

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
  console.log("server listening on port", port);
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
