done=make(chan bool)
on_message=fn(_script,_message,_data,_userdata){
  defer fn(){
    err=recover()
    if err!=nil{
      msg(err)
    }
  }()
  defer fn(){
    err=recover()
    if err!=nil{
      msg(err)
    }
  }()
  if len(_message)==0{
    return
  }
  if(_message["type"]=="send"){
    payload=_message["payload"]
    if(payload["done"]!=undefined){
      done<-true
      return
    }
    msg(payload)
  }elif(_message["type"]=="error"){
    msg(_message["description"])
    msg(_message["stack"])
  }elif(_message["type"]=="log"){
    msg(_message["payload"])
  }else{
    msg(sprintf("msg:%s",_message))
  }
}

openapp_with_d_bundleid=fn(_d,_bundleid){
  pid=0
  session=nil
  displayname=""
  apps,err=_d.enumApplication()
  if err!=nil{
    return 0,nil,"",err
  }
  for _,app=range apps{
    if app.identifier==_bundleid{
      pid=app.pid
      displayname=app.name
      if pid!=0{
        _d.kill(pid)  //有bug，先kill
        pid=0
      }
    }
  }

  if pid!=0{
    session,err=_d.attach(pid)
    if err!=nil{
      return 0,nil,"",err
    }
  }else{
    pid,err=_d.spawn(_bundleid)
    if err!=nil{
      return 0,nil,"",err
    }
    session,err=_d.attach(pid)
    if err!=nil{
      return 0,nil,"",err
    }
  msg(sprintf("进程id:%d 名称:%s",pid,displayname))
  script,err=session.Create_Script_with_path("./agent.js")
  if err!=nil{
    return err
  }
  script.on("message",on_message)
  err=script.load()
  if err!=nil{
    return err
  }
    err=_d.resume(pid)
    if err!=nil{
      return 0,nil,"",err
    }
  }
  return pid,session,displayname,nil
}
cmd_with_serial_cmd = fn (_serial,kwargs){
  defer fn(){
    err=recover()
    if err!=nil{
      msg(err)
    }
  }()

  manager=frida.fridadevicemanager()
  dev,err=manager.GetDevice_with_id_milltimeout(_serial,1000)
  if err!=nil{
    return err
  }

  defer manager.close()
  identifier=kwargs["identifier"]
  pid,session,displayname,err=openapp_with_d_bundleid(dev,identifier)
  if err!=nil{
    return err
  }


  //defer session.detach()
  <-done
  return nil
}
main{
  defer fn(){
    err=recover()
    if err!=nil{
      msg(err)
    }
  }()
  err=cmd_with_serial_cmd(args[0],{"identifier":"com.wondertek.hecmccmobile"})
  if err!=nil{
    msg(err)
  }
}
