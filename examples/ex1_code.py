

"""
import dist_comp as dc

dc.connect(
    ws_host="host", 
    pip=["numpy"]
)

# @dc.remote(numgpu=1)
def img_process(img_id:int):
    print(img_id)

a = [img_process.remote(i) for i in range(1,10)]


"""

from typing import Callable
import inspect as ins 
import ast 
from astunparse import unparse

def get_clean_source(func):
    source = ins.getsource(func)
    tree = ast.parse(source)

    # Find the first FunctionDef node
    for node in ast.walk(tree):
        if isinstance(node, ast.FunctionDef):
            node.decorator_list = []
            return unparse(node)
    raise ValueError("No FunctionDef found in source")



def remote_exec(gpu_required: bool):
    
    def remote_exec_wrap(func: Callable):
        source = get_clean_source(func)
        print("executing code remotely...")
        exec(source)

    return remote_exec_wrap

# def send_executable_code(code: )

# TODO: figure out what to pass params as.
@remote_exec(gpu_required=True)
def process_img(img_id:int):

    print(f"processing img {img_id}")


