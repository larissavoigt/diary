window.addEventListener("push", function(e) {
  createRatingGraph();
});

document.body.addEventListener("input", function(e) {
  var el = e.target;
  var ranges = {
    5: 'Desperately Unhappy',
    15: 'Very Unhappy',
    25: 'Unhappy',
    35: 'Passable',
    45: 'Quite OK',
    55: 'OK',
    65: 'Content',
    75: 'Cheerful',
    85: 'Happy',
    95: 'Delighted',
    100: 'Blissful'
  };
  if (el && el.id == "rate") {
    var name = document.getElementById("rate-label");
    var val = parseInt(el.value, 10);
    for (var i in ranges) {
      if (val <= i) {
        name.innerHTML = ranges[i];
        break;
      }
    }
  }
});

function createRatingGraph() {
  var chart = document.getElementById("graph");
  if (!chart) return;
  var ctx = chart.getContext("2d");
  var graph = new Chart(ctx).Bar({
    labels: [0,10,20,30,40,50,60,70,80,90,100],
    datasets: [
    {
      label: "My First dataset",
      fillColor: "rgba(220,220,220,0.5)",
      strokeColor: "rgba(220,220,220,0.8)",
      highlightFill: "rgba(220,220,220,0.75)",
      highlightStroke: "rgba(220,220,220,1)",
      data: graphData
    }
    ]
  });

}

createRatingGraph();
