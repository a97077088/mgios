import *as ios from "frida-lib/ios"
// @ts-ignore
import MGStatic = ObjC.classes.MGStatic;
// @ts-ignore
import MGAnalitics = ObjC.classes.MGAnalitics;
// @ts-ignore
import SessionAppInfo = ObjC.classes.SessionAppInfo;

// @ts-ignore
import SessionSDK=ObjC.classes.SessionSDK

try{
    // ios.fast.setExceptionHandle(null)
    setImmediate(function () {
        main()
    })

}catch(e){
    console.log(e)
}

function main(){
    // ios.fast.showcallmethod_with_hookoccls(SessionSDK)
    // ios.fast.showcallmethod_with_hookoccls(MGStatic)
    // ios.fast.showcallmethod_with_hookoccls(SessionAppInfo)
}
