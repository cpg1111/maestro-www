'use strict';

export default {
    handleDebugErr: function handleDebugErr(){
        if(window.DEBUG){
            arguments.forEach(console.error);
        }
    }
};
