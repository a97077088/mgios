base=0x102C181D2

for i in range(0,100):
    s=GetString(base,-1,0)
    print("\"{}\",".format(s))
    base+=len(s)+1