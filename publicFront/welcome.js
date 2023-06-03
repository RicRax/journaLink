   const query = window.location.search.substring(1);
   const token = query.split("access_token=")[1];

    fetch("https://api.github.com/user", {
      headers: {
        "Accept" : "application/vnd.github.+json",

        "Authorization": "Bearer " + token,

        "X-GitHub-Api-Version": "2022-11-28"
      },
    })
      .then((res) => res.json())
      .then((res) => {
        var welcome = document.getElementById('welcome');
        welcome.innerHTML = JSON.stringify(res.login);
      });

