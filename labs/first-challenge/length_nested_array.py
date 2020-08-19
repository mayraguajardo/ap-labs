
array = [1,2,3,[1,3],1,[5,6,7]]

def recursive_len(item):

    if type(item) == list:
        return sum(recursive_len(subitem)for subitem in item)
    else:
        return 1 

print(recursive_len(array))
