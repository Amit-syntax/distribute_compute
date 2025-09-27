"""
Author: amitkr20020405@gmail.com(amit-syntax)
"""

from dataclasses import dataclass
from typing import Callable
import inspect as ins
import ast
from astunparse import unparse
import websocket
import json


host = ""
port = ""
_ws_conn: websocket.WebSocket
_session_id = None


def connect(
    ws_host: str, ws_port: str,
    client_id: str, pip_modules: list[str],
):

    global host, port, _ws_conn, _session_id

    host = ws_host
    port = ws_port

    _ws_conn = websocket.create_connection(f"ws://{host}:{port}/session")
    _ws_conn.send(json.dumps({
        "pip_modules": pip_modules, "client_id": client_id,
    }))

    response = _ws_conn.recv()
    if type(response) is bytes:
        response = response.decode("utf-8")

    print("Response from server:", response)
    response = json.loads(response)

    _session_id = response["body"]["session_id"]
    print(f"Connected to server. Session ID: {_session_id}")


def get_clean_source(func):
    source = ins.getsource(func)
    tree = ast.parse(source)

    # Find the first FunctionDef node
    for node in ast.walk(tree):
        if isinstance(node, ast.FunctionDef):
            node.decorator_list = []
            return unparse(node), node.name
    raise ValueError("No FunctionDef found in source")


@dataclass
class RemExecObj:
    func_source: str
    func_args: str
    remote_call_uuid: str
    gpu_required: bool = True

    def __repr__(self) -> str:
        return f"RemExecObj({self.remote_call_uuid})"


class RemoteExec:

    def __init__(self, gpu_required: bool):
        self.gpu_required = gpu_required
        self.func = None

    def __call__(self, func: Callable):
        self.func = func
        return self

    def remote(self, *args, **kwargs) -> RemExecObj:

        args = str(args)
        source, func_name = get_clean_source(self.func)
        source += f"\n{func_name}" + args

        execution_uuid = self._remote_exec(
            source, func_name, args
        )

        return RemExecObj(
            source, args,
            execution_uuid,
            self.gpu_required,
        )

    def _remote_exec(
        self,
        func_source: str, func_name: str,
        func_args: str
    ) -> str:
        # TODO: send code to server and get result
        # - catch exceptions from server if can't execute
        rem_exec_call_obj = {
            "type": "session_remote_exec_request",
            "description": "Request to execute function remotely",
            "body": {
                "session_id": _session_id,
                "func_source": func_source,
                "func_args": func_args,
                "func_name": func_name,
                "params": {
                    "gpu_required": self.gpu_required,
                }
            }
        }
        _ws_conn.send(json.dumps(rem_exec_call_obj))

        resp = _ws_conn.recv()
        if type(resp) is bytes:
            resp = resp.decode("utf-8")
        resp = json.loads(resp)

        print("Response from server:", resp)
        execution_uuid = resp.get("body", {}).get("execution_id", None)

        return execution_uuid
