class GamePadMapper {
    constructor() {
        this.gamepadIndex = -1;
        this.numButtons = 32;
        this.numAxes = 10;
        this.bitButtons = 0;
        this.axes = Array(this.numAxes).fill(0.0);

        window.addEventListener('gamepadconnected', (event) => {
            const myGamepads = navigator.getGamepads();
            if (myGamepads != null && myGamepads[event.gamepad.index] != null) {
                this.gamepadIndex = event.gamepad.index;
            } else {
                console.log("Got event from null gamepad: ", event.gamepad.index);
            }
        });

        window.addEventListener('gamepaddisconnected', (evnet) => {
            this.gamepadIndex = -1;
        });

    }

    getGamePad() {
        if (this.gamepadIndex !== -1) {
            const myGamePads = navigator.getGamepads();
            const myGamePad = myGamePads[this.gamepadIndex];

            if (myGamePad.id.toLowerCase().includes("xbox") || myGamePad.id.toLowerCase().includes("0b20")) {
                return myGamePad;
            }
            // else if (myGamePad.id.toLowerCase().includes("g27")) {
            //     return myGamePad;
            // } else if (myGamePad.id.toLowerCase().includes("b684")) { //TGT wheel
            //     return myGamePad;
            // }
        }
        return null;
    }

    syncState(myGamePad) {
        if (myGamePad != null) {
            if (myGamePad.id.toLowerCase().includes("xbox") || myGamePad.id.toLowerCase().includes("0b20")) {
                return this.syncWithXbox(myGamePad);
            }
        }
        return false;
    }

    syncWithXbox(myGamePad) {
        this.bitButtons = 0; //reset state

        //spread across full range
        let mappedThrottle = this.mapToRange(myGamePad.buttons[7].value, 0.0, 1.0, -1.0, 1.0);
        let mappedBrake = this.mapToRange(myGamePad.buttons[6].value, 0.0, 1.0, -1.0, 1.0);


        this.axes[0] = Math.round(myGamePad.axes[0] * 1000) / 1000; //Steering
        this.axes[1] = Math.round(mappedThrottle * 1000) / 1000; //Throttle
        this.axes[2] = Math.round(mappedBrake * 1000) / 1000; //Brake
        this.axes[3] = Math.round(myGamePad.axes[2] * 1000) / 1000; //Pan
        this.axes[4] = Math.round(myGamePad.axes[3] * 1000) / 1000; //Tilt
        this.axes[5] = Math.round(myGamePad.axes[1] * 1000) / 1000; //Unused

        this.setBit(0, myGamePad.buttons[14].pressed);//steering left trim
        this.setBit(1, myGamePad.buttons[15].pressed);//steering right trim
        this.setBit(2, myGamePad.buttons[11].pressed);//camera recenter
        this.setBit(3, myGamePad.buttons[5].pressed);//upshift
        this.setBit(4, myGamePad.buttons[4].pressed);//downshift
        this.setBit(20, myGamePad.buttons[9].pressed);//client volume mute
        this.setBit(21, myGamePad.buttons[12].pressed);//client volume up
        this.setBit(22, myGamePad.buttons[13].pressed);//client volume down
        this.setBit(23, myGamePad.buttons[0].pressed);//unused
        this.setBit(24, myGamePad.buttons[1].pressed);//unused
        this.setBit(25, myGamePad.buttons[2].pressed);//unused
        this.setBit(26, myGamePad.buttons[3].pressed);//unused
        this.setBit(27, myGamePad.buttons[10].pressed);//unused
        return true;
    }

    getState() {
        return {
            "axes": this.axes,
            "bit_buttons": this.bitButtons,
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
        return Math.floor((maxReturn - minReturn) * (value - min) / (max - min) + minReturn)
    }
}