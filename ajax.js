  window.onload = function() {
    document.onReady
    document.getElementById("submit").addEventListener("click", submit);
  }

  makeblob = function (dataURL) {
            var BASE64_MARKER = ';base64,';
            if (dataURL.indexOf(BASE64_MARKER) == -1) {
                var parts = dataURL.split(',');
                var contentType = parts[0].split(':')[1];
                var raw = decodeURIComponent(parts[1]);
                return new Blob([raw], { type: contentType });
            }
            var parts = dataURL.split(BASE64_MARKER);
            var contentType = parts[0].split(':')[1];
            var raw = window.atob(parts[1]);
            var rawLength = raw.length;

            var uInt8Array = new Uint8Array(rawLength);

            for (var i = 0; i < rawLength; ++i) {
                uInt8Array[i] = raw.charCodeAt(i);
            }

            return new Blob([uInt8Array], { type: contentType });
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
    console.log("trying to poop");
    url = "http://derekspace.ddns.net:8080/save/";
    cFunction = function() {
      document.getElementById("success").innerHTML = "success!";
    };

    loadDoc(url, cFunction)
  }
