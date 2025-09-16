

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




