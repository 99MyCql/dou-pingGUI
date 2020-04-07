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
          <el-button size="small" type="primary">ping</el-button>
          <el-button size="small" type="info">stop</el-button>
        </el-form-item>
      </el-form>

      <el-progress type="circle" :percentage="progress"></el-progress>
    </div>

    <el-divider>结果</el-divider>

    <div id="body">
      <el-table
        :data="tableData"
        style="width: 100%"
        fit
        :row-class-name="tableRowClassName">
        <el-table-column
          type="index"
          label="#"
          align="center"
          width="100">
        </el-table-column>
        <el-table-column
          prop="date_len"
          label="返回字节 / Byte"
          align="center">
        </el-table-column>
        <el-table-column
          prop="time"
          label="时间 / ms"
          align="center">
        </el-table-column>
        <el-table-column
          prop="ttl"
          label="TTL"
          align="center">
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: {
        ip: '',
        count: 5,
        data: 'hello world!'
      },
      progress: 25,
      tableData: [{
        suc: true,
        date_len: 23,
        time: 12,
        ttl: 32,
      }, {
        suc: true,
        date_len: 23,
        time: 12,
        ttl: 32,
      }, {
        suc: true,
        date_len: 23,
        time: 12,
        ttl: 32,
      }, {
        suc: true,
        date_len: 23,
        time: 12,
        ttl: 32,
      }]
    }
  },
  methods: {
    tableRowClassName({rowIndex}) {
      if (this.tableData[rowIndex].suc) {
        return 'success-row';
      } else {
        return 'warning-row';
      }
    },
    countHandle(currentValue) {
      if (currentValue > 10) this.form.count = 10;
      else if (currentValue < 0) this.form.count = 0;
    }
  },
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

</style>