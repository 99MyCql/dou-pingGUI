<template>
  <div id="ipnet-ping">
    <div id="header">
      <el-form :inline="true" :model="form">
        <el-form-item label="目标子网" style="margin-bottom:0">
          <el-select
            style="width: 120px;margin-right: 10px;"
            size="small"
            v-model="form.type"
            @change="typeChange"
            placeholder="请选择">
            <el-option label="当前子网" value="cur_ipnet"></el-option>
            <el-option label="任意子网" value="arb_ipnet"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-show="this.form.type == 'cur_ipnet'" label="主机 IP" style="margin-bottom:0">
          <el-select
            style="width: 160px;margin-right: 30px;"
            size="small"
            v-model="form.ipNet"
            placeholder="请选择">
            <el-option v-for="(ip, i) in this.ipList" :key="i" :label="ip" :value="ip"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-show="this.form.type == 'arb_ipnet'" label="网络号(IP/CIDR)" style="margin-bottom:0">
          <el-input
            style="width: 160px;margin-right: 30px;"
            size="small"
            v-model="form.ipNet"
            placeholder="如：192.168.0.100/24">
          </el-input>
        </el-form-item>
        <el-form-item style="margin-bottom:0">
          <el-button size="small" type="primary" @click="ping">ping</el-button>
          <el-button size="small" type="danger" @click="stop">stop</el-button>
        </el-form-item>
      </el-form>

      <el-progress type="circle" :percentage="progress"></el-progress>
    </div>

    <el-divider>结果</el-divider>

    <div id="body">
      <el-table
        :data="tableData"
        :row-class-name="tableRowClassName"
        style="width:100%"
        height="360"
        :fit="true">
        <el-table-column
          type="index"
          label="#"
          align="center"
          width="100">
        </el-table-column>
        <el-table-column
          prop="ip"
          sortable
          :sort-method="ipCmp"
          label="目标 IP"
          align="center"
          min-width="191">
        </el-table-column>
        <el-table-column
          prop="time"
          sortable
          :sort-method="timeCmp"
          label="时间 / ms"
          align="center"
          min-width="191">
        </el-table-column>
        <el-table-column
          prop="ttl"
          sortable
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
        type: "cur_ipnet",
        ipNet: '',
      },
      counter: 0,   // 计数器，用于 tableData 中 id 赋值
      ipNetSize: 0,
      ipList: [],
      tableData: []
    }
  },
  computed: {
    progress: function() {
      if (this.tableData.length > 0) {
        return Math.round((this.tableData.length / this.ipNetSize) * 100);
      } else {
        return 0;
      }
    }
  },
  methods: {
    // time 排序规则：a<b return -1; a>b return 1; a==b return 0;
    timeCmp(a, b) {
      window.console.log(a.time, b.time);
      if ((a.time == "timeout" && b.time != "timeout")
        || parseInt(a.time) > parseInt(b.time)) {
        return 1;
      } else if ((a.time != "timeout" && b.time == "timeout")
        || parseInt(a.time) < parseInt(b.time)) {
        return -1;
      } else if ((a.time == "timeout" && b.time == "timeout")
        || parseInt(a.time) == parseInt(b.time)) {
        return 0;
      }
    },
    // ip 排序规则
    ipCmp(a, b) {
      if (a.ip.length < b.ip.length) {
        return -1;
      } else if (a.ip.length > b.ip.length) {
        return 1;
      } else {
        if (a.ip < b.ip) {
          return -1;
        } else if (a.ip > b.ip) {
          return 1;
        } else {
          return 0;
        }
      }
    },
    // 子网类型变化
    typeChange() {
      this.form.ipNet = '';
    },
    // 行的样式
    tableRowClassName({row}) {
      // window.console.log(row, rowIndex);
      if (this.tableData[row.id].suc) {
        return 'success-row';
      } else {
        return 'warning-row';
      }
    },
    // ping 按钮点击后调用
    ping() {
      this.counter = 0;
      this.tableData = [];
      window.backend.Controller.PingIPNet(this.form.ipNet)
      .then(ipNetSize => {
        if (ipNetSize <= 0) {
          this.$message.error('网络号（IP/CIDR）错误！');
        } else {
          this.ipNetSize = ipNetSize
        }
      })
    },
    // stop 按钮点击后调用
    stop() {
      window.backend.Controller.SetStop();
    }
  },
  mounted() {
    window.backend.Controller.GetIPAddrs()
    .then(data => {
      window.console.log(data);
      this.ipList = data;
    })

    Wails.Events.On("pingIPNet", (suc, ip, ttl, duration, data_len) => {
      window.console.log(suc, ip, ttl, duration, data_len)
      if (!suc) {
        duration = "timeout"
      }
      this.tableData.push({
        id: this.counter,
        suc: suc,
        ip: ip,
        data_len: data_len,
        time: duration,
        ttl: ttl,
      })
      this.counter++;
    });
  }
}
</script>

<style>
#ipnet-ping {
  padding: 0 80px 60px 80px;
}

#ipnet-ping #header {
  padding: 20px 0px;
  display: flex;
  justify-content: space-around;
  align-items: center;
}

.el-table .warning-row {
  background-color: #fef0f0;
  color: #f56c6c;
}

.el-table .success-row {
  background-color: #f0f9eb;
  color: #67c23a;
}
</style>