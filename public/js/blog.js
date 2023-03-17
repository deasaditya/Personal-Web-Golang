// global array
let blogs = []

// function get data from input text html
function getData(e) {
    e.preventDefault()

    // delaclaration variable dom selection
    let title = document.getElementById('create-project-name').value
    let startdate = document.getElementById('start-date').value
    let enddate = document.getElementById('end-date').value
    let content = document.getElementById("create-project-textarea").value;
    let nodejs = document.getElementById("node-js").checked ? document.getElementById("node-js").value : false;
    let reactjs = document.getElementById("react-js").checked ? document.getElementById("react-js").value : false;
    let nextjs = document.getElementById("next-js").checked ? document.getElementById("next-js").value : false;
    let typescript = document.getElementById("type-script").checked ? document.getElementById("type-script").value : false;
    let image = document.getElementById("image").files;

    if (title == "") 
  alert("harap isi nama");
  // else if (startdate == "") 
  // alert("silahkan tentukan tanggal mulai");
  // else if (enddate == "") 
  // alert("silahkan tentukan tanggal akhir");
  else if (content == "") 
  alert("isi deskripsi");
  else if (nodejs == false && reactjs == false && nextjs == false && typescript == false) 
  alert("silahkan centang salah satu technologies");
  else if 
  (image.length == 0) alert("pilih gambar terlebih dahulu");
  else { 
    image = URL.createObjectURL(image[0]);
    nodejs = nodejs != false ? `<i class="fa-brands ${nodejs} fa-3x"></i>` : false;
    reactjs = reactjs != false ? `<i class="fa-brands ${reactjs} fa-3x"></i>` : false;
    nextjs = nextjs != false ? `<i class="fa-brands ${nextjs} fa-3x"></i>` : false;
    typescript = typescript != false ? `<i class="fa-brands ${typescript} fa-3x"></i>` : false;
    

    let dataBlog = {
        title,
        content,
        nodejs,
        reactjs,
        nextjs,
        typescript,
        image,
        author: "Deas Aditya",
        postedAt: new Date()    
    };
    blogs.push(dataBlog);
    // console.log(blogs);
    showData();

    }
}

// declaration function show list preview data blog
function showData() {
  console.log(blogs);
  document.getElementById("containers-project").innerHTML = "";
  for (let blog of blogs) {
    document.getElementById("containers-project").innerHTML += `
    <div id="container-project">
      <img src="${blog.image}"/>
      <a href="detail.html"><h3>${blog.title}</h3></a>
      <div class="durasi">
      ${blog.postedAt} | ${blog.author}</div>
      <p>
        ${blog.content}
      </p>
      <div>
        ${blog.nodejs != false ? blog.nodejs : ""}
        ${blog.reactjs != false ? blog.reactjs : ""}
        ${blog.nextjs != false ? blog.nextjs : ""}
        ${blog.typescript != false ? blog.typescript : ""}
      </div>
      <div>
        <button>edit</button>
        <button>delete</button>
      </div>
      <div style="float:right; margin: 10px">
                <p style="font-size: 15px; color:grey">${getDuration(blog.postedAt)}</p>
            </div>
    </div>
    `;
  }
}


// function manipulation date time
function createTime(time){
    // declaration variable
    let years = time.getFullYear()
    let monthIndex = time.getMonth()
    let date = time.getDate()
    let hour = time.getHours()
    let minutes = time.getMinutes()
    
    const month = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]

    return `${date} ${month[monthIndex]} ${years} ${hour}:${minutes} WIB`
}


// manipulation duration time
const getDuration = (timePost) => {
    let timePosting = new Date(timePost)
    let timeNow = new Date(timePost)
    let distance = timeNow - timePosting
console.log(distance)
    // 1 sec = 1000 milisecond
    // 30*24*60*60*1000

    let monthDistance = Math.floor(distance / (30*24*60*60*1000))
    if(monthDistance > 0) {
        return monthDistance + ' month ago'
    } else {
        let weekDistance = Math.floor(distance / (7*24*60*60*1000))
        if(weekDistance > 0) {
            return weekDistance + ' week ago'
        } else {
            let dayDistance = Math.floor(distance / (24*60*60*1000))
            if(dayDistance > 0) {
                return dayDistance + "Day ago"
            } else {
                let hourDistance = Math.floor(distance / (60*60*1000))
                if(hourDistance > 0) {
                    return hourDistance + "Hour ago"
                } else {
                    let minutesDistance = Math.floor(distance / (60*1000))
                    if(minutesDistance > 0) {
                        return minutesDistance + "Minute ago"
                    } else {
                        let secondDistance = Math.floor(distance / 1000)
                        if (secondDistance > 0) {
                            return secondDistance + "second ago"
                        }
                    }
                }
            }
        }
    }
}

setInterval(showData, 1000)


