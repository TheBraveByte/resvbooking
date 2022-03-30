function promptAlert() {
    let toast = function (notify) {
        const {
            msg = "",
            icon = "success",
            position = "top",
        } = notify

        //Default parameter
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmationBar: false,
            timer: 5000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener("mouseenter", Swal.stopTimer)
                toast.addEventListener("mouseleave", Swal.resumeTimer)
            }
        })
        Toast.fire({})
    }

    let success = function (notify) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = notify

        //Default parameter
        Swal.fire({
            icon: "success",
            title: title,
            footer: footer
        })

    }

    let error = function (notify) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = notify

        //Default parameter
        Swal.fire({
            icon: "error",
            title: title,
            footer: footer
        })

    }


    async function customer(notify) {
        const {
            msg = "",
            title = "",
            icon = "",
            showConfirmationButton = true,

        } = notify

        const {
            value: formValues
        } = await Swal.fire({
            icon: icon,
            title: title,
            backdrop: false,
            showCancelButton: true,
            showConfirmationButton: showConfirmationButton,
            html: msg,
            focusConfirm: false,
            willOpen: () => {
                if (notify.willOpen !== undefined) {
                    notify.willOpen();
                }
            },

            preConfirm: () => {
                return [
                    document.getElementById("check-in").value,
                    document.getElementById("check-out").value
                ]
            },
            didOpen: () => {
                if (notify.didOpen() !== undefined) {
                    notify.didOpen();
                }

            }
        },)
        if (formValues) {
            if (formValues.dismiss !== Swal.DismissReason.cancel) {
                if (formValues !== "") {
                    if (notify.callback !== undefined) {
                        notify.callback(formValues)
                    } else {
                        notify.callback(false);
                    }
                }
            } else {
                notify.callback(false);
            }
        }

    }

    return {
        success: success,
        toast: toast,
        error: error,
        customer: customer,
    }

}