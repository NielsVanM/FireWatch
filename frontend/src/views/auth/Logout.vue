<template>
    <span>Logut View</span>
</template>

<script>
export default {
    name: "LogoutView",
    created: function() {
        this.$http
            .post(this.$store.state.backendURL+"/api/v1/logout/", {
                token: localStorage.getItem("token")
                })
            .then((resp) => {
                if (resp.data.success == true) {
                    this.$router.push("/login/")
                }
            })
    }
}
</script>

