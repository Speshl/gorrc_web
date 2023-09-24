class KeyPressMapper {
    constructor(){
        this.pressedKeys = {};

        this.maxPosition = 1.0;
        this.midPosition = 0;
        this.minPosition = -1.0;

        this.numButtons = 32;
        this.numAxes = 10;
        this.bitButtons = 0;
        this.axes = Array(this.numAxes).fill(0.0);

        // Event listener for keydown event
        document.addEventListener('keydown', (event) => {
            const key = event.key;
            this.pressedKeys[key] = true;
            if(key == "ArrowUp" || key == "ArrowDown" || key == "ArrowLeft" || key == "ArrowRight" || key == " "){
                event.preventDefault();
            }
        });

        // Event listener for keyup event
        document.addEventListener('keyup', (event) => {
            const key = event.key;
            delete this.pressedKeys[key];
        });
    }

    syncState() {
        this.bitButtons = 0; //reset state

        if(this.pressedKeys['a'] === true) { //steering
            this.axes[0] = this.minPosition;
        }else if(this.pressedKeys['d'] === true) {
            this.axes[0] = this.maxPosition;
        }else{
            this.axes[0] = this.midPosition;
        }

        if(this.pressedKeys['w'] === true) { //throttle
            this.axes[1] = this.maxPosition;
        }else{
            this.axes[1] = this.midPosition;
        }

        if(this.pressedKeys['s'] === true) { //brake
            this.axes[2] = this.maxPosition;
        }else{
            this.axes[2] = this.midPosition;
        }

        if(this.pressedKeys['ArrowLeft'] === true) { //pan
            this.axes[3] = this.minPosition;
        }else if(this.pressedKeys['ArrowRight'] === true) {
            this.axes[3] = this.maxPosition;
        }else{
            this.axes[3] = this.midPosition;
        }

        if(this.pressedKeys['ArrowDown'] === true) { //tilt
            this.axes[4] = this.minPosition;
        }else if(this.pressedKeys['ArrowUp'] === true) {
            this.axes[4] = this.maxPosition;
        }else{
            this.axes[4] = this.midPosition;
        }


        this.setBit(0, this.pressedKeys[',']);//steering left trim
        this.setBit(1, this.pressedKeys['.']);//steering right trim
        this.setBit(2, this.pressedKeys[' ']);//camera recenter
        this.setBit(3, this.pressedKeys['e']);//upshift
        this.setBit(4, this.pressedKeys['q']);//downshift
        this.setBit(20, this.pressedKeys['m']);//client volume mute
        this.setBit(21, this.pressedKeys[']']);//client volume up
        this.setBit(22, this.pressedKeys['[']);//client volume down
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
}