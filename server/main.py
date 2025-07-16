import asyncio
import websockets


async def handle_connection(websocket):
    
    async for message in websocket:
        print(f"Received message: {message}")
        response = f"this is server responce about your message: {message}"
        await websocket.send(response)
        print(f"Sent response: {response}")


async def main():
    host = "localhost"
    port = 8765

    # start the WebSocket server
    server = await websockets.serve(handle_connection, host, port)
    print(f"WebSocket server started at ws://{host}:{port}")
    await server.wait_closed()


if __name__ == "__main__":
    asyncio.run(main())
