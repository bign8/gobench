var GB = GB || {};

GB.run = function init() {
  if (GB.hasOwnProperty('points')) GB.plot();
};

GB.plot = (function plotWrap() {
  function unpack(rows, key) {
    return rows.map(function(row) { return row[key]; });
  }

  function pad(num) {
    if (num < 10) num = "0" + num;
    return num;
  }

  // convert to yyyy-mm-dd HH:MM:SS.ssssss format
  function fmt(stamp) {
    var s = new Date(stamp);
    var YMD = s.getFullYear() + "-" + pad(s.getMonth()) + "-" + pad(s.getDate());
    var HMS = pad(s.getHours()) + ":" + pad(s.getMinutes()) + ":" + pad(s.getSeconds());
    return YMD + " " + HMS;
  }

  return function plot() {
    var data = GB.points;
    var stamps = unpack(data, "stamp").map(fmt);

    function set(name) {
      return [{
        x: stamps,
        y: unpack(data, name),
        // type: 'scatter',
      }];
    }

    Plotly.plot("benchN", set("iter"), {
      title: "Number of Iterations",
    });
    Plotly.plot("benchNS", set("ns"), {
      title: "Nanoseconds per Operation",
    });
    Plotly.plot("benchB", set("b"), {
      title: "Bytes per Operation",
    });
    Plotly.plot("benchAlloc", set("allocs"), {
      title: "Allocations per Operation",
    });
  };
})();

// Initializer
// TODO: use https://gist.github.com/mrdoob/838785 if necessary
if (window.requestAnimationFrame) {
  requestAnimationFrame(GB.run);
} else {
  setTimeout(GB.run, 16);
}
