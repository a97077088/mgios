import * as ios from "frida-lib/ios"
import base = Kernel.base;
// @ts-ignore
var NSUbiquitousKeyValueStore=ObjC.classes.NSUbiquitousKeyValueStore

var UIDevice=ObjC.classes.UIDevice
// @ts-ignore
import ASIdentifierManager=ObjC.classes.ASIdentifierManager

import NSUserDefaults=ObjC.classes.NSUserDefaults

import NSString=ObjC.classes.NSString

import NSDictionary=ObjC.classes.NSDictionary

import NSObject=ObjC.classes.NSObject


function strToHexCharCode(str:any) {
    if(str === "")
        return "";
    var hexCharCode = [];
    hexCharCode.push("0x");
    for(var i = 0; i < str.length; i++) {
        hexCharCode.push((str.charCodeAt(i)).toString(16));
    }
    return hexCharCode.join("");
}
ios.fast.mainfunc(function () {////////////
    // @ts-ignore
    // var idfa=UIDevice.currentDevice().currentUserSessionId()//

    try{


        // @ts-ignore
        console.log(ObjC.classes.MGUser)
        // @ts-ignore
        var setProvinceCode=ObjC.classes.MGUser["- setProvinceCode:"]
        console.log(setProvinceCode)
        Interceptor.attach(setProvinceCode.implementation,{
            onEnter:function(args){
                console.log(`setProvinceCode:${new ObjC.Object(args[2])}`)
                ios.fast.show_backtrace()
            }
        })


        // var sha256=ios.libc.getExportFunction("CC_SHA256","void",["pointer","int","pointer"])
        // Interceptor.attach(sha256,{
        //     onEnter:function(args){
        //         // @ts-ignore
        //         console.log(`sha256:${hexdump(args[0],{length:Number(args[1]),header:false})}`)
        //         console.log()
        //     }
        // })
        //
        // var argx:NativePointer=NULL
        // var CCHMac=ios.libc.getExportFunction("CCHmac","void",["int","pointer","int","pointer","int","pointer"])
        // Interceptor.attach(CCHMac,{
        //     onEnter:function(args){
        //         // @ts-ignore
        //         console.log(`cchmac:${hexdump(args[1],{length:Number(args[2]),header:false})}`)
        //         console.log(`cchmac:${hexdump(args[3],{length:Number(args[4]),header:false})}`)
        //         console.log()
        //         ios.fast.show_backtrace()
        //         argx=args[5]
        //     },
        //     onLeave:function(r){
        //         console.log(`cchmac r:${hexdump(argx,{length:32,header:false})}`)
        //     }
        // })
        //
        //
        // var baseaddr:NativePointer=(<Module>Process.findModuleByName("MiguVideo")).base
        // Interceptor.attach(baseaddr.add(0x36998C),{
        //     onEnter:function(args){
        //         console.log(`0x3E1B820:${baseaddr.add(0x3E1B820).readU32()}`)
        //         console.log(`0x3E1B818:${baseaddr.add(0x3E1B818).readPointer().readCString(baseaddr.add(0x3E1B820).readU32())}`)
        //
        //         console.log(`0x36998c :${hexdump(args[0],{header:false,length:16})}`)
        //         console.log(`0x36998c :${args[1].readUtf8String()}`)
        //         console.log(`0x36998c :${args[2].readCString()}`)
        //         console.log(`0x36998c :${args[3].readCString()}`)
        //         console.log(`0x36998c :${hexdump(args[4],{length:32,header:false})}`)
        //         console.log(`0x36998c :${args[5]}`)
        //     }
        // })
        //
        //
        // Interceptor.attach(baseaddr.add(0x369680),{
        //     onEnter:function (args) {
        //         console.log(`369680参数:${args[0].readCString()}`)
        //         console.log(`369680参数:${args[1].readCString()}`)
        //         console.log(`369680参数:${args[2].readCString()}`)
        //         console.log(`369680参数:${args[3].readCString()}`)
        //     }
        // })
        //
        // Interceptor.attach(baseaddr.add(0x181E74),{//
        //     onEnter:function (args) {
        //         console.log(hexdump(args[2],{length:Number(args[3]),header:false}))
        //         console.log(`248B28参数:${args[0].readCString()}`)
        //         console.log(`248B28参数:${args[1]}`)
        //         console.log(`248B28参数:${args[2].readUtf8String(Number(args[3]))}`)
        //         console.log(`248B28参数:${args[3]}`)
        //     },onLeave:function(r){
        //         console.log(`248B28参数:${r.readCString()}`)
        //     }
        // })
        //
        // //@ts-ignore
        // var getNativeMac=ObjC.classes.MGUASDKKDF["+ getNativeMac:data:"]
        // Interceptor.attach(getNativeMac.implementation,{
        //     onEnter:function(args){
        //         ios.fast.show_backtrace()
        //         console.log(`${ObjC.selectorAsString(args[1])} ${((args[2]).readCString(16))} , ${new ObjC.Object(args[3])}`)
        //     },
        //     onLeave:function(r){
        //         console.log(`getNativeMac:${new ObjC.Object(r)}`)
        //     }
        // })
        //
        // //@ts-ignore
        // var charToNSHex=ObjC.classes.MGUASDKFormatTransform["+ charToNSHex:length:"]
        // Interceptor.attach(charToNSHex.implementation,{
        //     onEnter:function(args){
        //         console.log(`${ObjC.selectorAsString(args[1])} ${((args[2]).readCString())} , ${args[3]}`)
        //     },
        //     onLeave:function(r){
        //         console.log(`charToNSHex:${new ObjC.Object(r)}`)
        //     }
        // })
        // // // @ts-ignore//
        // // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.MGUASDKKDF)
        // //@ts-ignore
        // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.MGUnionAuthManager)
        // // @ts-ignore
        // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.MGSignInManager)////////////////
        //
        // // @ts-ignore
        // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.IDMPSafeKSMgmt)
        // // @ts-ignore
        // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.IDMPToken)
        // // @ts-ignore
        // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.IDMPParseParament)
        // // @ts-ignore
        // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.MGUnionAuthManager)
        // // @ts-ignore
        // //ios.fast.showcallmethod_with_hookoccls(ObjC.classes.MGUASDKFormatTransform)
        // // @ts-ignore
        // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.IDMPMD5)
        // // @ts-ignore
        // ios.fast.showcallmethod_with_hookoccls(ObjC.classes.IDMPUPMode)

        // // @ts-ignore
        // var charToNSHex=ObjC.classes.MGUASDKFormatTransform["+ charToNSHex:length:"]
        // Interceptor.attach(charToNSHex.implementation,{
        //     onEnter:function(args){
        //         console.log(`${ObjC.selectorAsString(args[1])} ${(args[2]).readCString()} , ${args[3]}`)
        //     },
        //     onLeave:function(r){
        //         console.log(`tonshex:${new ObjC.Object(r)}`)
        //     }
        // })
        //
        //
        // // @ts-ignore
        // var getMd5_32Bit_String=ObjC.classes.IDMPMD5["+ getMd5_32Bit_String:"]
        // Interceptor.attach(getMd5_32Bit_String.implementation,{
        //     onEnter:function(args){
        //         console.log(`${ObjC.selectorAsString(args[1])} ${new ObjC.Object(args[2])}`)
        //     },
        //     onLeave:function(r){
        //         console.log(`md5r:${new ObjC.Object(r)}`)
        //     }
        // })

        // @ts-ignore
        // var setobjforkey=ObjC.classes.NSMutableDictionary["- hookSetObject:forKey:"]
        // Interceptor.attach(setobjforkey.implementation,{
        //     onEnter:function(args){
        //         var obj=new ObjC.Object(args[3])
        //         if (obj.toString()=="KS"){
        //             console.log(ObjC.selectorAsString(args[1]))
        //             console.log(new ObjC.Object(args[2]))
        //             console.log(new ObjC.Object(args[3]))
        //             ios.fast.show_backtrace()
        //         }
        //     }
        // })



        // @ts-ignore
        // var getTokenWithUserName=ObjC.classes.IDMPToken["+ getTokenWithUserName:appId:uuid:andUser:"]
        // Interceptor.attach(getTokenWithUserName.implementation,{
        //     onEnter:function(args){
        //         console.log(ObjC.selectorAsString(args[1]))
        //         console.log(new ObjC.Object(args[2]))
        //         console.log(new ObjC.Object(args[3]))
        //         console.log(new ObjC.Object(args[4]))
        //         console.log(new ObjC.Object(args[5]))
        //         ios.fast.show_backtrace()
        //     }
        // })

        // @ts-ignore
        // var mod=ObjC.classes.IDMPParseParament["+ parseParamentFrom:"]
        // Interceptor.attach(mod.implementation,{
        //     onEnter:(args => {
        //         console.log(ObjC.selectorAsString(args[1]))
        //         console.log(new ObjC.Object(args[2]))
        //     }),
        //     onLeave:function(r){
        //         console.log(new ObjC.Object(r))
        //     }
        // })
        // var baseaddr:NativePointer=(<Module>Process.findModuleByName("MiguVideo")).base
        // Interceptor.attach(baseaddr.add(0x1DE514),{
        //     onEnter:function (args) {
        //         console.log(`rsa1参数:${new ObjC.Object(args[0])}`)
        //         console.log(`rsa1参数:${new ObjC.Object(args[1])}`)
        //         ios.fast.show_backtrace()
        //     }
        // })
        // @ts-ignore
        // var mod=ObjC.classes.MGUnionAuthManager["- getAccessTokenByConditionWithLoginType:username:content:countryCode:successBlock:failBlock:"]
        // Interceptor.attach(mod.implementation,{
        //     onEnter:function(args){
        //         console.log(new ObjC.Object(args[5]))
        //         ios.fast.show_backtrace()
        //     }
        // })

        //
        // setTimeout(function(){
        //     send({"done":true})
        // },1000)
        //
        //     // // @ts-ignore
        // console.log(ObjC.classes.MGSignInManager["- handleSDKLoginSuccessWithParam:"])
        // // @ts-ignore
        // ios.fast.hookoc(ObjC.classes.MGSignInManager["- handleSDKLoginSuccessWithParam:"],function(oldfn,self,sel,a1,){
        //
        //     console.log(`arg1:${new ObjC.Object(a1)}`)
        //     var r=oldfn(self,sel,a1,)
        //     console.log(`返回:${new ObjC.Object(r)}`)
        //     //return r;
        // })
        //
        //
        //     // // @ts-ignore
        //     // ios.fast.hookoc(ObjC.classes.IDMPAES128["+ AESEncryptWithKey:andData:"],function(oldfn,self,sel,arg1,arg2){
        //     //     console.log(new ObjC.Object(arg1))
        //     //     console.log(new ObjC.Object(arg2))
        //     //     var r=oldfn(self,sel,arg1,arg2)
        //     //     return r;
        //     // })
        //
        //     // @ts-ignore
        //     // ios.fast.hookoc(ObjC.classes.IDMPAES128["+ base64EncodingWithData:"],function(oldfn,self,sel,arg1){
        //     //     console.log(new ObjC.Object(arg1))
        //     //     var tm:ObjC.classes.NSString=<any>NSString.alloc()
        //     //     var ts=tm.initWithData_encoding_(<any>new ObjC.Object(arg1),4)
        //     //     console.log(`base64参数:${ts}`)
        //     //
        //     //     var r=oldfn(self,sel,arg1)
        //     //     console.log(`base64 r:${new ObjC.Object(r)}`)
        //     //     return r;
        //     // })
        //
        //     var baseaddr:NativePointer=(<Module>Process.findModuleByName("MiguVideo")).base
        //     // Interceptor.attach(baseaddr.add(0x23ED00),{
        //     //     onEnter:function (args) {
        //     //         console.log(`rsa参数:${new ObjC.Object(args[0])}`)
        //     //     }
        //     // })
        //     // Interceptor.attach(baseaddr.add(0x313F70),{
        //     //     onEnter:function(args){
        //     //         console.log(`rsa参数:${new ObjC.Object(args[0])}`)
        //     //     }
        //     // })
        //     // @ts-ignore
        //     // ios.fast.hookoc(ObjC.classes.IDMPMD5["+ getMd5_32Bit_String:"],function(oldfn,self,sel,a1){
        //     //     var r=oldfn(self,sel,a1)
        //     //     console.log(`md5 参数:${new ObjC.Object(a1)}`)
        //     //     console.log(`md5 返回:${new ObjC.Object(r)}`)
        //     //     return r
        //     // })
        //
        //     // // @ts-ignore
        //     // ios.fast.hookoc(ObjC.classes.NSString["- md5String"],function(oldfn,self,sel){
        //     //     var r=oldfn(self,sel)
        //     //     console.log(`md5 参数:${new ObjC.Object(self)}`)
        //     //     console.log(`md5 返回:${new ObjC.Object(r)}`)
        //     //     return r
        //     // })
        // @ts-ignore
        // ios.fast.hookoc(ObjC.classes.IDMPRSA_Encrypt_Decrypt["+ addPublicKey:"],function(oldfn,self,sel,a1){
        //     var r=oldfn(self,sel,a1)
        //     console.log(new ObjC.Object(a1))
        //     console.log(r)
        //     // ios.fast.show_backtrace()
        //     return r
        // })
        //
        //     // send({'done':true})
    }catch(e){
        console.log(e)
    }

    // send({"done":true})
})


