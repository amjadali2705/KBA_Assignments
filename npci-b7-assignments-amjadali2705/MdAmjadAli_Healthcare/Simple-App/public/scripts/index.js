const addData = async (event) => {
    event.preventDefault();

    const patientId = document.getElementById("patientId").value;
    const name = document.getElementById("name").value;
    const age = document.getElementById("age").value;
    const medicalHistory = document.getElementById("medicalHistory").value;

    const patientData = {
        patientId: patientId,
        name: name,
        age: age,
        medicalHistory: medicalHistory,
    }

    if (
        patientId.length == 0 ||
        name.length == 0 ||
        age.length == 0 ||
        medicalHistory.length == 0
    ) {
        alert("Please enter the data properly.");
    } else {
        try {
            const response = await fetch("/api/patient", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(patientData)
            })

            return alert("Patient Added");
        } catch (err) {
            alert("Error");
            console.log(err);
        }
    }
}

const readData = async (event) => {
    event.preventDefault();
    const patientId = document.getElementById("patientIdInput").value;

    if (patientId.length == 0) {
        alert("Please enter a valid ID.");
    } else {
        try {
            const response = await fetch(`/api/patient/${patientId}`);
            let responseData = await response.json();
            console.log("response data", responseData);
            alert(JSON.stringify(responseData));
        } catch (err) {
            alert("Error");
            console.log(err);
        }
    }
};