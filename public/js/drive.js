// Toggle dark mode
function toggleDarkMode() {
    const body = document.querySelector('body');
    body.classList.toggle('dark-mode',true);
    body.classList.toggle('light-mode',false);
}

function toggleLightMode() {
    const body = document.querySelector('body');
    body.classList.toggle('dark-mode',false);
    body.classList.toggle('light-mode',true);
}

function selectTrack(el) {
    if (el.className.indexOf("selectedTrack") >= 0) {
        el.className = el.className.replace("selectedTrack","");
        document.getElementById("selectTrackButton").disabled = true;
    }
    else {
        const rows = document.getElementsByClassName("selectedTrack");
        for(let i=0; i<rows.length; i++){
            rows[i].className = el.className.replace("selectedTrack","");
        }
        el.className  += "selectedTrack";
        htmx.trigger("#"+el.id, "trackSelected");
    }
}

function selectCar(el) {
    if (el.className.indexOf("selectedCar") >= 0) {
        el.className = el.className.replace("selectedCar","");
        document.getElementById("selectCarButton").disabled = true;
    }
    else {
        const rows = document.getElementsByClassName("selectedCar");
        for(let i=0; i<rows.length; i++){
            rows[i].className = el.className.replace("selectedCar","");
        }
        el.className  += "selectedCar";
        htmx.trigger("#"+el.id, "carSelected");
    }
}

function startConnecting(trackName, carName, seatNumber){
    const camPlayer = new CamPlayer();
    
    setTimeout(() => {
        camPlayer.startMicrophone().then(() => {
            camPlayer.sendConnect(trackName, carName, seatNumber);
        });
        //camPlayer.sendConnect(trackName, carName, seatNumber);
    }, 1000);
    

    const gamePadMapper = new GamePadMapper();
    const keyPressMapper = new KeyPressMapper();
    
    //Start listener loop for input commands
    var controllerName = "";
    setInterval(() => {
        let gamePad = gamePadMapper.getGamePad();

        let state = null;
        let newControllerName = "";
        if(gamePad != null){
            newControllerName = gamePadMapper.syncState(gamePad);
            state = gamePadMapper.getState();
           
        }else{
            newControllerName = keyPressMapper.syncState();
            state = keyPressMapper.getState();
        }

        if(newControllerName !== controllerName){
            controllerName = newControllerName;
            const controllerType = document.getElementById('controllerType');
            controllerType.innerHTML = controllerName;
        }

        if (camPlayer.gotRemoteDescription() && state !== null) {
            camPlayer.sendState(state);
        }
    }, 10); //Update at 100hz
}

//Setup htmx event listeners
document.body.addEventListener("startconnecting", (e) =>{
    startConnecting(e.detail.track, e.detail.car, e.detail.seat);
})

//Startup all the processes we need
toggleDarkMode();
