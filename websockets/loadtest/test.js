import { WebSocket } from "k6/experimental/websockets";

export const options = {
  vus: 1,
  iterations: 100,
};

export default function () {
  testBroadcastHandler();
}

function testEchoHandler() {
  const numMessages = 10;
  const ws = new WebSocket("ws://localhost:8080");

  ws.onopen = () => {
    for (let i = 0; i < numMessages; i++) {
      const echoMessage = {
        type: "echo",
        payload: `Hello Echo from K6! - Message ${i + 1}`,
      };
      ws.send(JSON.stringify(echoMessage));
      console.log(`Sent Echo: ${JSON.stringify(echoMessage)}`);
    }
    ws.close();
  };
}

function testBroadcastHandler() {
  const numMessages = 10;
  const ws = new WebSocket("ws://localhost:8080");

  ws.onopen = () => {
    for (let i = 0; i < numMessages; i++) {
      const broadcastMessage = {
        type: "broadcast",
        payload: `Hello Broadcast from K6! - Message ${i + 1}`,
      };
      ws.send(JSON.stringify(broadcastMessage));
      console.log(`Sent Broadcast: ${JSON.stringify(broadcastMessage)}`);
    }
    ws.close();
  };
}

// Run the test for the "echo" handler
export function echoTest() {
  testEchoHandler();
}

// Run the test for the "broadcast" handler
export function broadcastTest() {
  testBroadcastHandler();
}

// import { WebSocket } from "k6/experimental/websockets";

// export const options = {
//   vus: 10,
//   iterations: 10,
// };

// export default function () {
//   const numMessages = 10;
//   const ws = new WebSocket("ws://localhost:8080"); // Replace with your WebSocket server URL

//   ws.onopen = () => {
//     for (let i = 0; i < numMessages; i++) {
//       // Send "echo" message
//       // const echoMessage = { type: "echo", payload: "Hello Echo from K6!" };
//       // ws.send(JSON.stringify(echoMessage));
//       // console.log(`Sent Echo: ${JSON.stringify(echoMessage)}`);

//       //  Send "broadcast" message
//       const broadcastMessage = {
//         type: "broadcast",
//         payload: "Hello Broadcast from K6!",
//       };
//       ws.send(JSON.stringify(broadcastMessage));
//       console.log(`Sent Broadcast: ${JSON.stringify(broadcastMessage)}`);
//     }
//     ws.close();
//   };
// }
