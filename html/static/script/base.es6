// base
// @author: zhengxiaoyao0716
; (function () {
    "use strict";
    const base = {
        "fetch": window.fetch ? url => fetch(url).then(
            r => r.status >= 200 && r.status < 300 ? Promise.resolve(r.json()) : Promise.reject(new Error(r.statusText))
        ) : url => new Promise((resolve, reject) => {
            const xmlhttp = window.XMLHttpRequest ? new XMLHttpRequest() : new ActiveXObject("Microsoft.XMLHTTP");
            xmlhttp.open('Get', url, true);
            xmlhttp.responseType = 'json';
            xmlhttp.addEventListener("load", e => xmlhttp.status == 200 ? resolve(xmlhttp.response) : reject(xmlhttp.statusText), false);
            xmlhttp.addEventListener("error", e => reject(e), false);
            xmlhttp.send();
        }),
        "preload": {
            "image": (src, onLoaded, onFailed, onComplete) => {
                onLoaded = onLoaded || function () { };
                onFailed = onFailed || function () { };
                onComplete = onComplete || function () { };

                const image = new Image();
                image.src = src;
                image.onerror = function (e) {
                    onFailed.call(this, src, e);
                    onComplete.call(this, src);
                };
                if (image.complete) {
                    onLoaded.call(image, src);
                    onComplete.call(image, src);
                } else {
                    image.addEventListener("load", function (e) {
                        onLoaded.call(this, src, e)
                        onComplete.call(this, src);
                    });
                }
            }
        }
    };

    // Module defined.
    if (typeof module !== 'undefined' && module.exports) {
        module.exports = base;
    } else if (typeof define === 'function' && define.amd) {
        define(function () { return base; });
    } else {
        window.base = base;
    }

    HTMLCollection.prototype.forEach = HTMLCollection.prototype.forEach || Array.prototype.forEach;
    NodeList.prototype.forEach = NodeList.prototype.forEach || Array.prototype.forEach;
})();
