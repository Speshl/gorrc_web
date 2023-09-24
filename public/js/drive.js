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

function driveCar(el,trackName, carName) {
    htmx.trigger("#"+el.id, "driveCar");
    document.body.addEventListener("startconnecting", () =>{
        startConnecting(trackName, carName);
    })
}

function startConnecting(trackName, carName){
    const camPlayer = new CamPlayer();
    
    setTimeout(() => {
        // camPlayer.startMicrophone().then(() => {
        //     camPlayer.sendOffer();
        // });
        //camPlayer.sendOffer();
        camPlayer.sendConnect(trackName, carName);
    }, 1000);
    

    const gamePadMapper = new GamePadMapper();
    const keyPressMapper = new KeyPressMapper();
    
    //Start listener loop for input commands
    setInterval(() => {
        let gamePad = gamePadMapper.getGamePad();

        let state = null;
        if(gamePad != null){
            gamePadMapper.syncState(gamePad);
            state = gamePadMapper.getState();
        }else{
            keyPressMapper.syncState();
            state = keyPressMapper.getState();
        }

        if (camPlayer.gotRemoteDescription() && state !== null) {
            camPlayer.sendState(state);
        }
    }, 17); //Update at 60hz
}


//Startup all the processes we need
toggleDarkMode();
