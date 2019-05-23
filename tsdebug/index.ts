import *as ios from "frida-lib/ios"
// @ts-ignore
import MGStatic = ObjC.classes.MGStatic;
// @ts-ignore
import MGAnalitics = ObjC.classes.MGAnalitics;
// @ts-ignore
import SessionAppInfo = ObjC.classes.SessionAppInfo;

// @ts-ignore
import HBRSAHandler=ObjC.classes.HBRSAHandler

try{
    // ios.fast.setExceptionHandle(null)
    setImmediate(function () {
        main()
    })

}catch(e){
    console.log(e)
}

function main(){


    Interceptor.attach(HBRSAHandler["- signMD5String:"].implementation,{
        onEnter:function(args){
            console.log(new ObjC.Object(args[2]))
        }
    })
}
