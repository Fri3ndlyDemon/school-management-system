document.addEventListener("DOMContentLoaded", function () {
    const attendanceForm = document.querySelector("#attendance-form");

    if (attendanceForm) {
        attendanceForm.addEventListener("submit", function (e) {
            // Optional client-side validation or dynamic feedback
            const studentId = document.querySelector("#student_id").value;
            const date = document.querySelector("#date").value;
            
            if (!studentId || !date) {
                e.preventDefault();
                alert("Please fill in all required attendance fields.");
            }
        });
    }

    // Add interactive styling toggles for status buttons if needed
    const statusSelect = document.querySelector("#status");
    if (statusSelect) {
        statusSelect.addEventListener("change", function () {
            if (this.value === "Present") {
                this.style.borderColor = "#28a745";
            } else {
                this.style.borderColor = "#dc3545";
            }
        });
    }
});