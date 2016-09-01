var GB = GB || {};

GB.run = function() {
  if (GB.hasOwnProperty('points')) GB.plot();
};

GB.plot = (function plotWrap() {
  function unpack(rows, key) {
    return rows.map(function(row) { return row[key]; });
  }

  return function plot() {
    var stamps = unpack(GB.points, "stamp").map(Plotly.d3.time.format.iso.parse);

    function set(name) {
      return [{
        x: stamps,
        y: unpack(GB.points, name),
        // type: 'scatter',
      }];
    }

    Plotly.plot("benchN", set("iter"), {
      title: "Number of Iterations",
    });
    GB.out = Plotly.plot("benchNS", set("ns"), {
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
window.requestAnimationFrame ? requestAnimationFrame(GB.run) : setTimeout(GB.run, 16);
