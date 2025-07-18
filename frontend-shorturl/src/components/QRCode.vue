<template>
  <div class="qrcode-container">
    <div ref="qrcode" class="qrcode"></div>
    <el-button 
      v-if="showDownload"
      size="mini" 
      @click="downloadQRCode">
      下载二维码
    </el-button>
  </div>
</template>

<script>
import QRCode from 'qrcodejs2'
export default {
  props: {
    text: {
      type: String,
      required: true
    },
    size: {
      type: Number,
      default: 120
    },
    showDownload: {
      type: Boolean,
      default: true
    }
  },
  mounted() {
    this.generate()
  },
  methods: {
    generate() {
      this.$refs.qrcode.innerHTML = ''
      new QRCode(this.$refs.qrcode, {
        text: this.text,
        width: this.size,
        height: this.size,
        colorDark: '#000000',
        colorLight: '#ffffff',
        correctLevel: QRCode.CorrectLevel.H
      })
    },
    downloadQRCode() {
      const canvas = this.$refs.qrcode.querySelector('canvas')
      const link = document.createElement('a')
      link.download = 'short-url-qrcode.png'
      link.href = canvas.toDataURL('image/png')
      link.click()
    }
  },
  watch: {
    text() {
      this.generate()
    }
  }
}
</script>

<style scoped>
.qrcode-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.qrcode {
  margin-bottom: 10px;
}
</style>    