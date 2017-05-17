// layer
// @author: zhengxiaoyao0716
; (function () {
    "use strict";
    // 图层随光标浮动
    ((containers) => {
        addEventListener('mousemove', e => containers.forEach((container) => {
            container.children.forEach((child) => {
                child.style.transformOrigin = 100 * e.x / container.clientWidth + '% ' + 100 * e.y / container.clientHeight + '%';
            });
        }));
        // addEventListener('click', e => containers.forEach((container) => {
        //     container.children.forEach((child) => {
        //         child.style.transformOrigin = 100 * (1 - e.x / container.clientWidth) + '% ' + 100 * (1 - e.y / container.clientHeight) + '%';
        //     });
        // }));
    })(document.getElementsByClassName('-layer-'));
})();
