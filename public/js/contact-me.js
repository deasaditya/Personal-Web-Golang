function getData(){
    let name = document.getElementById('specificSizeInputName').value
    let email = document.getElementById('exampleInputEmail1').value
    let phone = document.getElementById('InputNumber').value
    let subject = document.getElementById('subject').value
    let description = document.getElementById('floatingTextarea2').value

    if(name == "") {
        alert("Nama harus di isi")
    } else if(email == "") {
        alert("email harus di isi")
    } else if(phone == "") {
        alert("phone harus di isi")
    } else if(subject == "") {
        alert("subject harus di isi")
    } else if(description == "") {
        alert("description harus di isi")
    }

    const defaultEmail = "hi.adityadeas86@gmail.com"

    let mailTo = document.createElement('a')
    mailTo.href = `mailto:${defaultEmail}?subject=${subject}&body=Halo nama saya ${name}, saya ingin ${description} tolong hubungin saya ${phone}`
    mailTo.click()

    let audience = {
        name,
        email,
        phone,
        subject,
        description
    }

    console.log(audience)
}