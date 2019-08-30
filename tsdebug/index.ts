import * as ios from "frida-lib/ios"


ios.fast.mainfunc(function () {

    // @ts-ignore
    console.log(ObjC.classes.NSUserDefaults.standardUserDefaults().objectForKey_("MGTIDMPToken"))
})
