import *as ios from "frida-lib/ios"
// @ts-ignore
import MGStatic = ObjC.classes.MGStatic;

// @ts-ignore
import RequestData=ObjC.classes.RequestData
import NSString = ObjC.classes.NSString;
import NSData = ObjC.classes.NSData;
try{
    // ios.fast.setExceptionHandle(null)
    setImmediate(function () {
        main()
    })

}catch(e){
    console.log(e)
}

function main(){


    // //aquireToken
    Interceptor.attach(RequestData["+ requestVerifyWithURL:withMethod:withData:withResult:"].implementation,{
        onEnter:function(args){
            var a2=new ObjC.Object(args[2])
            var a4=(<NSString>NSString.alloc()).initWithData_encoding_(new ObjC.Object(args[4]),1)
            console.log(`aquireToken`)
            console.log()
        }
    })
    //dataCollectionService
    Interceptor.attach(RequestData["+ requestQulityWithURL:withMethod:withToken:withData:withResult:"].implementation,{
        onEnter:function(args){
            var a2=new ObjC.Object(args[2])
            var a5=(<NSString>NSString.alloc()).initWithData_encoding_(new ObjC.Object(args[5]),1)
            console.log(`dataCollectionService原始:${a5}`)
            var evs=events_with_s(a5)
            for(var ev of evs){
                console.log(ev)
            }
            console.log()
        }
    })
    //dataExposureService
    Interceptor.attach(RequestData["+ requestConfigWithURL:withMethod:withToken:withData:withResult:"].implementation,{
        onEnter:function(args){
            var a2=new ObjC.Object(args[2])
            var a5=(<NSString>NSString.alloc()).initWithData_encoding_(new ObjC.Object(args[5]),1)
            console.log(`dataExposureService原始:${a5}`)
            var evs=events_with_s(a5)
            for(var ev of evs){
                console.log(ev)
            }
            console.log()
        }
    })

    //dataEventService
    Interceptor.attach(RequestData["+ requestDefineWithURL:withMethod:withToken:withData:withResult:"].implementation,{
        onEnter:function(args){
            var a2=new ObjC.Object(args[2])
            var a5=(<NSString>NSString.alloc()).initWithData_encoding_(new ObjC.Object(args[5]),1)
            console.log(`dataEventService原始:${a5}`)
            var evs=events_with_s(a5)
            for(var ev of evs){
                console.log(ev)
            }
            console.log()
        }
    })
}

var tss={
    "4":"logIdentify",
    "5":"disConnectToServer",
    "8":"programDownloadFailed",
    "9":"playFailed",
    "10":"firstPageDelay",
    "11":"videoPageLoadDelay",
    "12":"searchPageLoadDelay",
    "14":"videoDownloadSpeed",
    "15":"heartBeatMessage",
    "16":"logIdentify",
    "17":"videoPlayStatus",
    "18":"purchaseStatistics",
    "19":"picloadDelay",
    "20":"pushMessageStatistics",
    "21":"exitApp",
    "22":"getPlayUrlDuration",
    "25":"downloadFlow",
    "26":"playFlow",
    "27":"crashexitdetail",
    "28":"errorLog",
    "37":"subSessionDownloadFlow",
    "26001000":"MGDownloadBytesEvent",
    "56100002":"MGPlayerGSLBReqEvent",
    "56100000":"MGPlayerCreateEvent",
    "56100001":"MGPlayerSetURLEvent",
    "56100009":"MGPlayerStartCmdEvent",
    "56100005":"MGPlayerM3U8ReqEvent",
    "56100006":"MGPlayerM3U8ParseEvent",
    "56100003":"MGPlayerDNSParseEvent",
    "56100004":"MGPlayerTCPOpenEvent",
    "56100010":"MGPlayerTSReqEvent",
    "56100014":"MGPlayerHttpConnectEvent",
    "56100007":"MGPlayerMediaOpenEvent",
    "56100008":"MGPlayerMediaInfoEvent",
    "56100011":"MGPlayerDecoderInitEvent",
    "56100012":"MGPlayerContiSeekEvent",
    "56100013":"MGPlayerFirstDecodeEvent",
    "56000015":"MGStuckEvent",
    "56000004":"MGFirstVideoRenderEvent",
    "57000000":"MGTrafficStatusticsEvent",
    "58000000":"MGErrorEvent",
    "56000016":"MGPlayerDecoderSwitchEvent",
    "59000000":"MGUserOperationEvent ",
    "70000000":"MGPlayerShutdownEvent",
    "60000000":"MGPlayerDurationEvent",
    "60100000":"MGPlayerAVDiffEvent",
    "80000000":"MGAuthenticationEvent",
}
function name_with_type(_s:string):string{
    var r=tss[_s]
    if(r==null){
        return "不支持解析的名字"
    }
    return r
}
function events_with_s(_s:string):Array<string>{
    var r=new Array<string>()
    try{
        var jso=JSON.parse(_s)
        if(jso.sessionStart!=null&&Object.keys(jso.sessionStart).length>0){
            r.push(`<sessionStart>->${JSON.stringify(jso.sessionStart)}`)
        }else if(jso.sessionEnd!=null&&Object.keys(jso.sessionEnd).length>0){
            r.push(`<sessionEnd>->${JSON.stringify(jso.sessionEnd)}`)
        }else if(jso.customInfo!=null&&Object.keys(jso.customInfo).length>0){
            for(var i=0;i<jso.customInfo.length;i++){
                r.push(`<${name_with_type(jso.customInfo[i].type)}:${jso.customInfo[i].type}>->${JSON.stringify(jso.customInfo[i])}`)
            }
        }else if(jso.customEvent!=null&&Object.keys(jso.customEvent).length){
            for(var i=0;i<jso.customEvent.length;i++){
                r.push(`<${jso.customEvent[i].eventName}>->${JSON.stringify(jso.customEvent[i])}`)
            }
        }else{
            console.log("不支持的解析格式")
        }
    }catch(e){
        console.log(e)
    }
    return r
}
