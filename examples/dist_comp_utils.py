import inspect as ins 
import ast 
from astunparse import unparse

def get_clean_source(func):
    source = ins.getsource(func)
    tree = ast.parse(source)
    # Find the first FunctionDef node
    for node in ast.walk(tree):
        if isinstance(node, ast.FunctionDef):
            node.decorator_list = []  # Remove decorators
            return unparse(node)
    raise ValueError("No FunctionDef found in source")

def remote(func):
    print("the source code of func: ")
    print(">",get_clean_source(func))

    return func

@remote
def some():
    print('this is a remote function')


some()
# remote(some)
# print(getattr(some, "__name__"))
# print(getattr(some, "__module__"))
# print(getattr(some, "__qualname__"))
# print(getattr(some, "__doc__"))
# print(getattr(some, "__annotations__"))

