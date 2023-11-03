class GamePadMapper {
    constructor() {
        this.connectedGamepadIndexes = [];
        this.preferredGamepadIndex = -1;
        this.numButtons = 32;
        this.numAxes = 10;
        this.bitButtons = 0;
        this.deadzone = 0.10;
        this.axes = Array(this.numAxes).fill(0.0);

        window.addEventListener('gamepadconnected', (event) => {
            console.log("gamepad connected event");
            const myGamepads = navigator.getGamepads();
            if (myGamepads != null && myGamepads[event.gamepad.index] != null) {

                let found = false;
                for(let i=0; i<this.connectedGamepadIndexes.length; i++){
                    if (this.connectedGamepadIndexes[i] === event.gamepad.index){
                        found = true;
                        break;
                    }
                }

                if(!found){
                    this.connectedGamepadIndexes.push(event.gamepad.index);
                    console.log("connected gamepad id " + event.gamepad.id + " @ index: " + event.gamepad.index);
                }
            } else {
                console.log("Got event from null gamepad: ", event.gamepad.index);
            }
        });

        window.addEventListener('gamepaddisconnected', (event) => {
            console.log("gamepad disconnected event");
            for(let i=0; i<this.connectedGamepadIndexes.length; i++){
                if (this.connectedGamepadIndexes[i] === event.gamepad.index){
                    console.log("disconnected gamepad id " + event.gamepad.id + " @ index: " + event.gamepad.index);
                    this.connectedGamepadIndexes.splice(i, 1); // 2nd parameter means remove one item only
                    break;
                }
            }
        });

    }

    getGamePad() {
        if (this.gamepadIndex !== -1) {
            const myGamePads = navigator.getGamepads();
            const myGamePad = myGamePads[this.gamepadIndex];
            if(this.isSupportedGamepad(myGamePad.id) != "unknown"){
                return myGamePad;
            }
        }
        return null;
    }

    getGamePad() {
        const myGamePads = navigator.getGamepads();
        if(this.preferredGamepadIndex !== -1 && myGamePads[this.preferredGamepadIndex] !== null){
            if(this.isSupportedGamepad(myGamePads[this.preferredGamepadIndex]) != "unknown"){
                return myGamePads[this.preferredGamepadIndex];
            }
        }

        for(let i=0; i<this.connectedGamepadIndexes.length; i++){
            let checkIndex = this.connectedGamepadIndexes[i];
            if(this.isSupportedGamepad(myGamePads[checkIndex].id) != "unknown"){
                return myGamePads[checkIndex];
            }
        }
        return null;
    }

    syncState(myGamePad) {
        if (myGamePad != null) {
            gamePadName = this.isSupportedGamepad(myGamePad.id);
            switch(gamePadName){
                case "xbox":
                    return this.syncWithXbox(myGamePad);
                case "dualshock":
                    return this.syncWithPlaystation(myGamePad);
                case "g27":
                    return this.syncWithG27(myGamePad);
            }
            return gamePadName;
        }
        return "unknown";
    }

    syncWithXbox(myGamePad) {
        this.bitButtons = 0; //reset state
        this.axes[0] = this.mapAxisWithDeadzone(myGamePad.axes[0], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Steering
        this.axes[1] = this.mapTriggerWithDeadzone(myGamePad.buttons[7].value, 0.0, 1.0, -1.0, 1.0, this.deadzone);//Throttle
        this.axes[2] = this.mapTriggerWithDeadzone(myGamePad.buttons[6].value, 0.0, 1.0, -1.0, 1.0, this.deadzone);//Brake
        this.axes[3] = this.mapAxisWithDeadzone(myGamePad.axes[2], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Pan
        this.axes[4] = this.mapAxisWithDeadzone(myGamePad.axes[3], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Tilt
        this.axes[5] = this.mapAxisWithDeadzone(myGamePad.axes[1], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Unused
        this.axes[9] = 0.3; //Used for steering sensitivy curve

        this.setBit(0, myGamePad.buttons[14].pressed);//steering left trim
        this.setBit(1, myGamePad.buttons[15].pressed);//steering right trim
        this.setBit(2, myGamePad.buttons[11].pressed);//camera recenter
        this.setBit(3, myGamePad.buttons[5].pressed);//upshift
        this.setBit(4, myGamePad.buttons[4].pressed);//downshift
        this.setBit(20, myGamePad.buttons[9].pressed);//client volume mute
        this.setBit(21, myGamePad.buttons[12].pressed);//client volume up
        this.setBit(22, myGamePad.buttons[13].pressed);//client volume down
        // this.setBit(23, myGamePad.buttons[0].pressed);//unused
        // this.setBit(24, myGamePad.buttons[1].pressed);//unused
        // this.setBit(25, myGamePad.buttons[2].pressed);//unused
        // this.setBit(26, myGamePad.buttons[3].pressed);//unused
        // this.setBit(27, myGamePad.buttons[10].pressed);//unused
        return "xbox";
    }

    syncWithPlaystation(myGamePad) {
        this.bitButtons = 0; //reset state
        this.axes[0] = this.mapAxisWithDeadzone(myGamePad.axes[0], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Steering
        this.axes[1] = this.mapTriggerWithDeadzone(myGamePad.buttons[7].value, 0.0, 1.0, -1.0, 1.0, this.deadzone);//Throttle
        this.axes[2] = this.mapTriggerWithDeadzone(myGamePad.buttons[6].value, 0.0, 1.0, -1.0, 1.0, this.deadzone);//Brake
        this.axes[3] = this.mapAxisWithDeadzone(myGamePad.axes[2], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Pan
        this.axes[4] = this.mapAxisWithDeadzone(myGamePad.axes[3], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Tilt
        this.axes[5] = this.mapAxisWithDeadzone(myGamePad.axes[1], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Unused
        this.axes[9] = 0.3; //Used for steering sensitivy curve

        this.setBit(0, myGamePad.buttons[14].pressed);//steering left trim
        this.setBit(1, myGamePad.buttons[15].pressed);//steering right trim
        this.setBit(2, myGamePad.buttons[11].pressed);//camera recenter
        this.setBit(3, myGamePad.buttons[5].pressed);//upshift
        this.setBit(4, myGamePad.buttons[4].pressed);//downshift
        this.setBit(20, myGamePad.buttons[9].pressed);//client volume mute
        this.setBit(21, myGamePad.buttons[12].pressed);//client volume up
        this.setBit(22, myGamePad.buttons[13].pressed);//client volume down
        // this.setBit(23, myGamePad.buttons[0].pressed);//unused
        // this.setBit(24, myGamePad.buttons[1].pressed);//unused
        // this.setBit(25, myGamePad.buttons[2].pressed);//unused
        // this.setBit(26, myGamePad.buttons[3].pressed);//unused
        // this.setBit(27, myGamePad.buttons[10].pressed);//unused
        return "dualshock";
    }

    syncWithG27(myGamePad) {
        this.bitButtons = 0; //reset state
        this.axes[0] = this.mapAxisWithDeadzone(myGamePad.axes[0], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Steering
        this.axes[1] = this.mapToRange(myGamePad.axes[6]*-1, -1.0, 1.0, -1.0, 1.0);//Throttle (inverted)
        this.axes[2] = this.mapToRange(myGamePad.axes[5]*-1, -1.0, 1.0, -1.0, 1.0);//Brake (inverted)
        this.axes[5] = this.mapToRange(myGamePad.axes[2], -1.0, 1.0, -1.0, 1.0, this.deadzone);//Unused (clutch)

        this.axes[9] = 0.0; //Used for steering sensitivy curve

        //axes 9 is dpad with set value for each button
        switch(myGamePad.axes[9].toFixed(2)){
            case -1.00: //UP
                this.axes[4] = 1.0;
                break;
            case -0.71: //UP & RIGHT
                this.axes[4] = 1.0;
                this.axes[3] = 1.0;
                break;
            case -0.43: //RIGHT
                this.axes[3] = 1.0;
                break;
            case -0.14: //DOWN & RIGHT
                this.axes[3] = 1.0;
                this.axes[4] = -1.0;
                break;
            case 0.14: //DOWN
                this.axes[4] = -1.0;
                break;
            case 0.43: //DOWN AND LEFT
                this.axes[4] = -1.0;
                this.axes[3] = -1.0;
                break;
            case 0.71: //LEFT
                this.axes[3] = -1.0;
                break;
            case 1.00: //UP & LEFT
                this.axes[4] = 1.0;
                this.axes[3] = -1.0;
                break;
        }

        this.setBit(0, myGamePad.buttons[7].pressed);//steering left trim
        this.setBit(1, myGamePad.buttons[6].pressed);//steering right trim
        this.setBit(2, myGamePad.buttons[21].pressed);//camera recenter
        this.setBit(3, myGamePad.buttons[4].pressed);//upshift
        this.setBit(4, myGamePad.buttons[5].pressed);//downshift
        this.setBit(5, myGamePad.buttons[0].pressed);//trans toggle (sequential, hpattern)
        this.setBit(6, myGamePad.buttons[14].pressed);//reverse gear
        this.setBit(7, myGamePad.buttons[8].pressed);//first gear
        this.setBit(8, myGamePad.buttons[9].pressed);//second gear
        this.setBit(9, myGamePad.buttons[10].pressed);//third gear
        this.setBit(10, myGamePad.buttons[11].pressed);//fourth gear
        this.setBit(11, myGamePad.buttons[12].pressed);//fifth gear
        this.setBit(12, myGamePad.buttons[13].pressed);//sixth gear
        this.setBit(20, myGamePad.buttons[19].pressed);//client volume mute
        this.setBit(21, myGamePad.buttons[20].pressed);//client volume up
        this.setBit(22, myGamePad.buttons[22].pressed);//client volume down
        return "g27";
    }

    getState() {
        return {
            "axes": this.axes,
            "bit_buttons": this.bitButtons,
            "time_stamp": Date.now(),
        }
    }

    setBit(position, value) {
        if (value) {
            this.bitButtons |= (1 << position); // Set the bit to 1
        } else {
            this.bitButtons &= ~(1 << position); // Set the bit to 0
        }
    }

    isBitSet(position) {
        return (this.bitButtons & (1 << position)) !== 0;
    }

    mapToRange(value, min, max, minReturn, maxReturn) {
        return (maxReturn - minReturn) * (value - min) / (max - min) + minReturn;
    }

    mapTriggerWithDeadzone(value, min, max, minReturn, maxReturn, deadzone) {
        if (value > deadzone){
            let maxWithDeadzone = max - deadzone;
            let valueWithDeadzone = value - deadzone;
            let newRange = this.mapToRange(valueWithDeadzone, min, maxWithDeadzone, minReturn, maxReturn);
            let rounded = Math.round(newRange*100) / 100;
            return rounded;
        }
         return minReturn;
    }

    mapAxisWithDeadzone(value, min, max, minReturn, maxReturn, deadzone) {
        if (value > deadzone || value < (deadzone * -1)){
            let valueWithDeadzone = value - deadzone;
            if (value < (deadzone * -1)){
                valueWithDeadzone = value + deadzone;
            }
            let maxWithDeadzone = max - deadzone;
            let minWithDeadzone = min + deadzone;
            let newRange = this.mapToRange(valueWithDeadzone, minWithDeadzone, maxWithDeadzone, minReturn, maxReturn);
            let rounded = Math.round(newRange*100) / 100;
            return rounded;
        }
        return 0.0; //return midpoint
    }

    isSupportedGamepad(gamePadId) {
        if (gamePadId != null) {
            if (gamePadId.toLowerCase().includes("xbox") || gamePadId.toLowerCase().includes("0b20") || gamePadId.toLowerCase().includes("0207") || gamePadId.toLowerCase().includes("1532")) {
                return "xbox";
            }else if(gamePadId.toLowerCase().includes("Dualshock") || gamePadId.toLowerCase().includes("054C") || gamePadId.toLowerCase().includes("09cc")) {
                return "dualshock";
            }else if (gamePadId.toLowerCase().includes("g27") || gamePadId.toLowerCase().includes("046d") || gamePadId.toLowerCase().includes("c29b")) {
                return "g27";
            }
        }
        return "unknown"
    }
}