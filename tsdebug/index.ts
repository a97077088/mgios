import *as ios from "frida-lib/ios"
// @ts-ignore
import MGStatic = ObjC.classes.MGStatic;
// @ts-ignore
import MGAnalitics = ObjC.classes.MGAnalitics;
// @ts-ignore
import SessionAppInfo = ObjC.classes.SessionAppInfo;

// @ts-ignore
import MGPlayerInstance=ObjC.classes.MGPlayerInstance

try{
    // ios.fast.setExceptionHandle(null)
    setImmediate(function () {
        main()
    })

}catch(e){
    console.log(e)
}

function main(){
    Interceptor.attach(MGPlayerInstance["- handleSQMEvent:"].implementation,{
        onEnter:function (args) {
            var a2=new ObjC.Object(args[2])
            console.log(a2.$className)
        }
    })
    // ios.fast.showcallmethod_with_hookoccls(MGPlayerInstance)
    // ios.fast.showcallmethod_with_hookoccls(SessionSDK)
    // ios.fast.showcallmethod_with_hookoccls(MGStatic)
    // ios.fast.showcallmethod_with_hookoccls(SessionAppInfo)
}
