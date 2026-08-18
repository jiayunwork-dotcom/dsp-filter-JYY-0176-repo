# dsp-filter

dsp-filter 是一个数字滤波器设计与频率响应工具。输入设计规格（FIR 或 IIR、阶数、归一化截止、窗类型），输出滤波器系数、单位圆上的幅频（dB）/相频响应，以及零极点与稳定性判定。

## 约定与边界

- **归一化截止 fc**：相对采样率 fs=1，范围 (0, 0.5)，Nyquist = 0.5。频率网格同样取值 [0, 0.5]。
- **FIR 阶数 N**：抽头数 = N+1。本仓钉死类型 I（偶对称线性相位），因此 FIR 阶数必须为偶数，奇数阶返回 error。
- **窗类型**：`rect` / `hann`（`hanning` 视为同一窗）/ `hamming`。未知窗名返回 error。
- **IIR**：Butterworth 低通 + 双线性变换，含截止预畸（Ωc = 2·tan(π·fc)）；`a[0]` 归一为 1。阶数限制 [1, 16]。
- **频率响应**：H(e^jω) 由 `b`/`a` 直接计算，输出幅度 dB、解缠后的相位（弧度）、群延时（样本）。空网格或含 NaN/越界的网格返回 error。
- **零极点**：实系数多项式用 Durand-Kerner 求根；任一极点模长 > 1+1e-6 判为不稳定并在结果中标记。

## 构建与运行

```bash
go build ./...
go run . -http :8080          # Web UI + API，浏览器打开 http://localhost:8080
go run . design -kind fir -order 30 -cutoff 0.2 -window hamming
go run . design -kind iir -order 4 -cutoff 0.2
go run . response example/fir_hamming.json
go run . zplane example/fir_hamming.json
```

Docker：

```bash
docker build -t dsp-filter:probe -f benzhi.Dockerfile .
docker run --rm dsp-filter:probe go test ./...
```

## API

- `POST /api/design`：`{"kind":"fir|iir","order":30,"cutoff":0.2,"window":"hamming"}` → `{"kind":"fir","b":[...],"a":[1]}`
- `POST /api/response`：`{"b":[...],"a":[...],"points":256}` 或 `{"b":...,"a":...,"freq":[...]}` → `{"freq":[...],"mag_db":[...],"phase":[...],"group_delay":[...]}`
- `POST /api/zplane`：`{"b":[...],"a":[...]}` → `{"zeros":[{re,im}...],"poles":[...],"stable":true}`

失败返回非 2xx 与机器可读错误：`{"code":"unknown_window","message":"..."}`。常见 code：`unknown_window`、`cutoff_out_of_range`、`bad_order`、`bad_kind`、`empty_frequency_grid`、`bad_frequency`。

## 算例

- `example/fir_hamming.json`：30 阶 Hamming 窗低通，fc=0.2。
- `example/fir_hann.json`：40 阶 Hann 窗低通，fc=0.15。
- `example/iir_bw4.json`：4 阶 Butterworth 低通，fc=0.2（-3 dB 点与 fc 对齐）。
- `example/iir_bw8.json`：8 阶 Butterworth 低通，fc=0.1。

## 包结构

- `internal/design`：规格校验、窗函数、FIR 窗函数法、Butterworth 模拟原型、双线性变换与预畸。
- `internal/response`：频率网格、H(e^jω) 计算、幅度 dB、相位解缠、群延时。
- `internal/zplane`：多项式求根（Durand-Kerner）、零极点、稳定性与极点分类。
- `internal/api`：HTTP 服务与错误 JSON。
