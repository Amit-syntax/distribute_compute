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
_ws_conn = None

def connect(
    ws_host: str, ws_port: str, 
    client_id: str, pip_modules: list,
):

    global host, port, _ws_conn

    host = ws_host
    port = ws_port

    _ws_conn = websocket.create_connection(f"ws://{host}:{port}/session")
    _ws_conn.send(json.dumps({
        "pip_modules": pip_modules, "client_id": client_id,
    }))
    response = _ws_conn.recv()


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
    gpu_required: bool = True

class RemoteExec:
    
    def __init__(self, gpu_required: bool):
        self.gpu_required = gpu_required
        self.func = None

    def __call__(self, func: Callable):
        self.func = func
        return self

    def remote(self, *args, **kwargs) -> RemExecObj:
        source, func_name = get_clean_source(self.func)
        source += f"\n{func_name}" + str(args)

        return RemExecObj(
            source, str(args),
            self.gpu_required,
        )

def get_results(rem_exec_objs: list[RemExecObj]):
    # TODO: send code to server and get result
    print("executing code...")

    results = []

    # TODO: here send functions to server and ask for execution
    # - catch exceptions from server if can't execute
    for rec in rem_exec_objs:
        results.append(exec(rec.func_source))
    
    print(results)



@RemoteExec(gpu_required=True)
def process_img(img_id: int):
    print(f"processing img {img_id}")


rem_exec_objs = [process_img.remote(i) for i in range(100)]

get_results(rem_exec_objs)
