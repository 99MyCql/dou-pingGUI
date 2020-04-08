<template>
  <div id="ip-ping">
    <div id="header">
      <el-form :inline="true" :model="form">
        <el-form-item label="主机" style="margin-bottom:0">
          <el-input
            style="width: 130px;"
            class="form-item"
            size="small"
            v-model="form.ip"
            placeholder="IP or Host name">
          </el-input>
        </el-form-item>
        <el-form-item label="数据" style="margin-bottom:0">
          <el-input
            style="width: 150px;"
            class="form-item"
            size="small"
            v-model="form.data"
            maxlength="16"
            show-word-limit
            placeholder="any data to send">
          </el-input>
        </el-form-item>
        <el-form-item label="次数" style="margin-bottom:0">
          <el-input-number
            style="margin-right:20px"
            size="small"
            v-model="form.count"
            :min="1"
            :max="10"
            :change="countHandle">
          </el-input-number>
        </el-form-item>
        <el-form-item style="margin-bottom:0">
          <el-button size="small" type="primary" @click="ping">ping</el-button>
          <el-button size="small" type="danger" @click="stop">stop</el-button>
        </el-form-item>
      </el-form>

      <el-progress type="circle" :percentage="getProgress()"></el-progress>
    </div>

    <el-divider>结果</el-divider>

    <div id="body">
      <el-table
        :data="tableData"
        :row-class-name="tableRowClassName"
        style="width:100%"
        :fit="true">
        <el-table-column
          type="index"
          label="#"
          align="center"
          width="100">
        </el-table-column>
        <el-table-column
          prop="ip"
          label="目标 IP"
          align="center"
          min-width="191">
        </el-table-column>
        <el-table-column
          prop="data_len"
          label="返回字节 / Byte"
          align="center"
          min-width="191">
        </el-table-column>
        <el-table-column
          prop="time"
          label="时间 / ms"
          align="center"
          min-width="191">
        </el-table-column>
        <el-table-column
          prop="ttl"
          label="TTL"
          align="center"
          min-width="191">
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
import Wails from '@wailsapp/runtime';

export default {
  data() {
    return {
      form: {
        ip: '',
        count: 4,
        data: 'hello world!'
      },
      tableData: []
    }
  },
  methods: {
    getProgress() {
      if (this.tableData.length > 0) {
        return (this.tableData.length / this.form.count).toFixed(2) * 100;
      } else {
        return 0;
      }
    },
    tableRowClassName({rowIndex}) {
      if (this.tableData[rowIndex].suc) {
        return 'success-row';
      } else {
        return 'warning-row';
      }
    },
    countHandle(currentValue) {
      this.tableData = [];
      if (currentValue > 10) this.form.count = 10;
      else if (currentValue < 0) this.form.count = 0;
    },
    ping() {
      this.tableData = []
      window.backend.Controller.Ping(this.form.ip, this.form.data, this.form.count)
      .then(suc => {
        if (!suc) {
          this.$message.error('host name error!');
        }
      })
    },
    stop() {
      window.backend.Controller.SetStop();
    }
  },
  mounted() {
    Wails.Events.On("ping", (suc, ip, ttl, duration, data_len) => {
      window.console.log(suc, ip, ttl, duration, data_len)
      if (!suc) {
        duration = "timeout"
      }
      this.tableData.push({
        suc: suc,
        ip: ip,
        data_len: data_len,
        time: duration,
        ttl: ttl,
      })
    });
  }
}
</script>

<style>
#ip-ping {
  padding: 0 80px 60px 80px;
}

#ip-ping #header {
  padding: 20px 0px;
  display: flex;
  justify-content: space-around;
  align-items: center;
}

.form-item {
  margin-right: 10px;
}

.el-table .warning-row {
  background-color: #fef0f0;
  color: #f56c6c;
}

.el-table .success-row {
  background-color: #f0f9eb;
  color: #67c23a;
}

.el-table .my-cell {
  width: 30%;
}
</style>