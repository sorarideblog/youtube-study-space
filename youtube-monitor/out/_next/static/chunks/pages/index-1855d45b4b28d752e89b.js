_N_E=(window.webpackJsonp_N_E=window.webpackJsonp_N_E||[]).push([[6],{"/EDR":function(t,e,n){(window.__NEXT_P=window.__NEXT_P||[]).push(["/",function(){return n("23aj")}])},"23aj":function(t,e,n){"use strict";n.r(e),n.d(e,"default",(function(){return H}));var o=n("nKUr");function r(t,e){if(!(t instanceof e))throw new TypeError("Cannot call a class as a function")}function i(t,e){for(var n=0;n<e.length;n++){var o=e[n];o.enumerable=o.enumerable||!1,o.configurable=!0,"value"in o&&(o.writable=!0),Object.defineProperty(t,o.key,o)}}function a(t,e,n){return e&&i(t.prototype,e),n&&i(t,n),t}function s(t){if(void 0===t)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return t}function c(t,e){return(c=Object.setPrototypeOf||function(t,e){return t.__proto__=e,t})(t,e)}function u(t,e){if("function"!==typeof e&&null!==e)throw new TypeError("Super expression must either be null or a function");t.prototype=Object.create(e&&e.prototype,{constructor:{value:t,writable:!0,configurable:!0}}),e&&c(t,e)}function f(t){return(f="function"===typeof Symbol&&"symbol"===typeof Symbol.iterator?function(t){return typeof t}:function(t){return t&&"function"===typeof Symbol&&t.constructor===Symbol&&t!==Symbol.prototype?"symbol":typeof t})(t)}function l(t,e){return!e||"object"!==f(e)&&"function"!==typeof e?s(t):e}function h(t){return(h=Object.setPrototypeOf?Object.getPrototypeOf:function(t){return t.__proto__||Object.getPrototypeOf(t)})(t)}var p=n("rePB"),d=n("q1tI"),_=n.n(d),m=n("BzQq"),y=n.n(m);function v(t){var e=function(){if("undefined"===typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"===typeof Proxy)return!0;try{return Date.prototype.toString.call(Reflect.construct(Date,[],(function(){}))),!0}catch(t){return!1}}();return function(){var n,o=h(t);if(e){var r=h(this).constructor;n=Reflect.construct(o,arguments,r)}else n=o.apply(this,arguments);return l(this,n)}}var g=function(t){u(n,t);var e=v(n);function n(t){var o;return r(this,n),o=e.call(this,t),Object(p.a)(s(o),"intervalId",void 0),o.state={now:new Date},o}return a(n,[{key:"componentDidMount",value:function(){var t=this;this.intervalId=setInterval((function(){t.setState({now:new Date})}),1e3)}},{key:"componentWillUnmount",value:function(){this.intervalId&&clearInterval(this.intervalId)}},{key:"render",value:function(){return Object(o.jsxs)("div",{id:y.a.clock,children:[Object(o.jsxs)("div",{className:y.a.dateString,children:[this.state.now.getMonth()+1," / ",this.state.now.getDate()," /"," ",this.state.now.getFullYear()]}),Object(o.jsxs)("div",{className:y.a.timeString,children:[this.state.now.getHours(),"\uff1a",this.state.now.getMinutes()<10?"0"+this.state.now.getMinutes().toString():this.state.now.getMinutes()]})]})}}]),n}(_.a.Component),j=n("EYK1"),b=n.n(j);function x(t){var e=function(){if("undefined"===typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"===typeof Proxy)return!0;try{return Date.prototype.toString.call(Reflect.construct(Date,[],(function(){}))),!0}catch(t){return!1}}();return function(){var n,o=h(t);if(e){var r=h(this).constructor;n=Reflect.construct(o,arguments,r)}else n=o.apply(this,arguments);return l(this,n)}}var w=function(t){u(n,t);var e=x(n);function n(){return r(this,n),e.apply(this,arguments)}return a(n,[{key:"render",value:function(){if(this.props.default_room_state&&this.props.no_seat_room_state){var t=this.props.default_room_state.seats.length+this.props.no_seat_room_state.seats.length;return Object(o.jsxs)("div",{id:b.a.message,children:["Currently ",t," people working! \ud83d\udd25"]})}return Object(o.jsx)("div",{id:b.a.message})}}]),n}(_.a.Component),R=n("Sas2"),O=n.n(R),S=n("4oeA"),D=n.n(S);function P(t){var e=function(){if("undefined"===typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"===typeof Proxy)return!0;try{return Date.prototype.toString.call(Reflect.construct(Date,[],(function(){}))),!0}catch(t){return!1}}();return function(){var n,o=h(t);if(e){var r=h(this).constructor;n=Reflect.construct(o,arguments,r)}else n=o.apply(this,arguments);return l(this,n)}}var k=function(t){u(n,t);var e=P(n);function n(){return r(this,n),e.apply(this,arguments)}return a(n,[{key:"seatWithSeatId",value:function(t,e){var n=e[0];return e.forEach((function(e){e.seat_id===t&&(n=e)})),n}},{key:"render",value:function(){var t=this;if(this.props.layout&&this.props.roomState){var e=this.props.roomState.seats.map((function(t){return t.seat_id})),n=this.props.layout,r={widthPx:900*n.room_shape.width/n.room_shape.height,heightPx:900},i=r.widthPx*n.font_size_ratio,a={width:100*n.seat_shape.width/n.room_shape.width,height:100*n.seat_shape.height/n.room_shape.height},s=n.seats.map((function(t){return{x:100*t.x/n.room_shape.width,y:100*t.y/n.room_shape.height}})),c=n.partitions.map((function(t){for(var e,o,r=n.partition_shapes,i=t.shape_type,a=0;a<r.length;a++)r[a].name===i&&(e=100*r[a].width/n.room_shape.width,o=100*r[a].height/n.room_shape.height);return{widthPercent:e,heightPercent:o}})),u=n.partitions.map((function(t){return{x:100*t.x/n.room_shape.width,y:100*t.y/n.room_shape.height}})),f=n.seats.map((function(n,r){var c=e.includes(n.id),u=e.includes(n.id)?t.seatWithSeatId(n.id,t.props.roomState.seats).user_display_name:"";return Object(o.jsxs)("div",{className:D.a.seat,style:{backgroundColor:c?"#ee989f":"#fce7d2",left:s[r].x+"%",top:s[r].y+"%",width:a.width+"%",height:a.height+"%",fontSize:i+"px"},children:[Object(o.jsx)("p",{className:D.a.seatNum,style:{fontWeight:"bold"},children:n.id}),Object(o.jsx)("p",{className:D.a.userDisplayName,style:{overflow:"hidden"},children:u.substr(0,4)})]},n.id)})),l=n.partitions.map((function(t,e){return Object(o.jsx)("div",{className:D.a.partition,style:{left:u[e].x+"%",top:u[e].y+"%",width:c[e].widthPercent+"%",height:c[e].heightPercent+"%"}},t.id)}));return Object(o.jsx)("div",{children:Object(o.jsxs)("div",{id:D.a.roomLayout,style:{width:r.widthPx+"px",height:r.heightPx+"px"},children:[f,l]})})}return Object(o.jsx)("div",{children:"Loading"})}}]),n}(_.a.Component);function N(t){var e=function(){if("undefined"===typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"===typeof Proxy)return!0;try{return Date.prototype.toString.call(Reflect.construct(Date,[],(function(){}))),!0}catch(t){return!1}}();return function(){var n,o=h(t);if(e){var r=h(this).constructor;n=Reflect.construct(o,arguments,r)}else n=o.apply(this,arguments);return l(this,n)}}var E=function(t){u(n,t);var e=N(n);function n(){return r(this,n),e.apply(this,arguments)}return a(n,[{key:"render",value:function(){return Object(o.jsx)("div",{id:O.a.defaultRoom,children:Object(o.jsx)(k,{layout:this.props.layout,roomState:this.props.roomState})})}}]),n}(_.a.Component),I=n("7N6H"),C=n.n(I);function L(t){var e=function(){if("undefined"===typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"===typeof Proxy)return!0;try{return Date.prototype.toString.call(Reflect.construct(Date,[],(function(){}))),!0}catch(t){return!1}}();return function(){var n,o=h(t);if(e){var r=h(this).constructor;n=Reflect.construct(o,arguments,r)}else n=o.apply(this,arguments);return l(this,n)}}var z=function(t){u(n,t);var e=L(n);function n(){return r(this,n),e.apply(this,arguments)}return a(n,[{key:"render",value:function(){if(this.props.no_seat_room_state){var t=this.props.no_seat_room_state.seats.length;return Object(o.jsxs)("div",{id:C.a.standingRoom,children:[Object(o.jsx)("h2",{children:"\u3000"}),Object(o.jsx)("h2",{children:"Standing Room"}),Object(o.jsxs)("h3",{children:["\uff08Enter with ",Object(o.jsx)("span",{className:C.a.commandString,children:"!0"}),"\uff09"]}),Object(o.jsxs)("h2",{children:[t," People"]})]})}return Object(o.jsx)("div",{id:C.a.standingRoom})}}]),n}(_.a.Component),M=function(t,e){return n=fetch(t,e),new Promise((function(t,e){n.then((function(n){n.ok?n.json().then((function(e){t(e)})).catch((function(t){e(t)})):e(n)})).catch((function(t){e(t)}))}));var n};function W(t){var e=function(){if("undefined"===typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"===typeof Proxy)return!0;try{return Date.prototype.toString.call(Reflect.construct(Date,[],(function(){}))),!0}catch(t){return!1}}();return function(){var n,o=h(t);if(e){var r=h(this).constructor;n=Reflect.construct(o,arguments,r)}else n=o.apply(this,arguments);return l(this,n)}}var H=function(t){u(n,t);var e=W(n);function n(t){var o;return r(this,n),o=e.call(this,t),Object(p.a)(s(o),"intervalId",void 0),o.state={layout:null,default_room_state:null,no_seat_room_state:null},o}return a(n,[{key:"componentDidMount",value:function(){var t=this;this.intervalId=setInterval((function(){M("https://taa4p9klha.execute-api.ap-northeast-1.amazonaws.com/rooms_state").then((function(e){e.default_room.seats.forEach((function(t){return console.log(t.seat_id,t.user_display_name)})),console.log("fetch\u5b8c\u4e86"),t.setState({layout:e.default_room_layout,default_room_state:e.default_room,no_seat_room_state:e.no_seat_room})})).catch((function(t){return console.error(t)}))}),1500)}},{key:"componentWillUnmount",value:function(){this.intervalId&&clearInterval(this.intervalId)}},{key:"render",value:function(){return Object(o.jsxs)("div",{style:{height:1080,width:1920,margin:0,position:"relative"},children:[Object(o.jsx)(g,{}),Object(o.jsx)(w,{default_room_state:this.state.default_room_state,no_seat_room_state:this.state.no_seat_room_state}),Object(o.jsx)(E,{layout:this.state.layout,roomState:this.state.default_room_state}),Object(o.jsx)(z,{no_seat_room_state:this.state.no_seat_room_state})]})}}]),n}(_.a.Component)},"4oeA":function(t,e,n){t.exports={roomLayout:"DefaultRoomLayout_roomLayout__3GEjH",seat:"DefaultRoomLayout_seat__11-WG",partition:"DefaultRoomLayout_partition__2utnI",seatNum:"DefaultRoomLayout_seatNum__3YkZ7",userDisplayName:"DefaultRoomLayout_userDisplayName__1wSzE"}},"7N6H":function(t,e,n){t.exports={standingRoom:"StandingRoom_standingRoom__2bRlz",commandString:"StandingRoom_commandString__1QJmT"}},BzQq:function(t,e,n){t.exports={clock:"Clock_clock__1VNI5",dateString:"Clock_dateString__1sJUD",timeString:"Clock_timeString__1ixRO"}},EYK1:function(t,e,n){t.exports={message:"Message_message__2Drwu"}},Qetd:function(t,e,n){"use strict";var o=Object.assign.bind(Object);t.exports=o,t.exports.default=t.exports},Sas2:function(t,e,n){t.exports={defaultRoom:"DefaultRoom_defaultRoom__aHA-O"}},rePB:function(t,e,n){"use strict";function o(t,e,n){return e in t?Object.defineProperty(t,e,{value:n,enumerable:!0,configurable:!0,writable:!0}):t[e]=n,t}n.d(e,"a",(function(){return o}))}},[["/EDR",0,1]]]);