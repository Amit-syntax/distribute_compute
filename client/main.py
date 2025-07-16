import asyncio
import websockets


async def run_my_code_loop(websocket: websockets.ClientConnection):
    
    while True:
        message = input("enter the message you want to send: ")
        if message.lower() == 'quit':
            break
            
        await websocket.send(message)
        resp_msg = await websocket.recv()
        print("msg recv:", resp_msg)


async def test_connect():

    host_name = input("Enter the server address:")
    port = input("Enter the port:")
    host = f"ws://{host_name}:{port}"

    try:
        async with websockets.connect(host) as ws:
            await run_my_code_loop(ws)
            
    except websockets.exceptions.ConnectionClosed:
        print("connection closed")
    except ConnectionRefusedError:
        print("couldn't connect to server")


if __name__ == "__main__":
    asyncio.run(test_connect())