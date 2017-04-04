  window.onload = function() {
    document.onReady
    document.getElementById("submit").addEventListener("click", submit);
  }

  function loadDoc(url, cFunction) {
    var xhttp = new XMLHttpRequest();
    xhttp.open("POST", url, true);
    //xhttp.setRequestHeader("Content-type", "application/json");
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4 && this.status == 200) {
        //var json = JSON.parse(xhttp.responseText);
        //console.log(json.email + ", " + json.password)
        cFunction(this);
      }
    };
    title = document.getElementById("title").value;
    description = document.getElementById("description").value;
    image = document.getElementById("image").value;
    video = document.getElementById("video").value;
    tags = document.getElementById("tags").value;
    categories = document.getElementById("categories").value;
    draft = document.getElementById("draft").value;
    body = document.getElementById("body").value;
    var data = JSON.stringify({
      "title": title,
      "description": description,
      "image": image,
      "video": video,
      "tags": tags,
      "categories": categories,
      "draft": draft,
      "body": body
    });
    xhttp.send(data);
  }

  function submit() {
    url = "http://localhost:8080/save/";
    cFunction = function() {
      document.getElementById("success").innerHTML = "success!";
    };

    loadDoc(url, cFunction)
  }
