on getMicrophoneVolume()
  input volume of (get volume settings)
end getMicrophoneVolume

on disableMicrophone()
  set volume input volume 0
end disableMicrophone

on enableMicrophone()
  set volume input volume 100
end enableMicrophone

if getMicrophoneVolume() is greater than 0 then
  disableMicrophone()
  return "0"
else
  enableMicrophone()
  return "100"
end if