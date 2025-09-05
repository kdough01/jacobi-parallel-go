import subprocess
import matplotlib.pyplot as plt
import numpy as np

rows_per_task_options = [1, 2, 3, 4]
num_threads = [2, 4, 6, 8, 12]
desired_convergence = "0.000001"
max_iter = "1000"
matrix_size = "1024"
num_runs = 5

def get_average_runtime(mode, threads=None, rows_per_task=None):
    total_time = 0.0
    for i in range(num_runs):
        cmd = [
            "go", "run", "benchmark.go",
            f"-mode={mode}",
            f"-convergence={desired_convergence}",
            f"-maxIter={max_iter}",
            f"-size={matrix_size}"
        ]

        if threads is not None:
            cmd.append(f"-numThreads={threads}")
        if rows_per_task is not None:
            cmd.append(f"-rowsPerTask={rows_per_task}")

        print(f"Run {i+1}/{num_runs}: {' '.join(cmd)}")
        result = subprocess.run(cmd, capture_output=True, text=True)

        try:
            output = result.stdout.strip().splitlines()
            for line in reversed(output):
                try:
                    run_time = float(line)
                    total_time += run_time
                    break
                except ValueError:
                    continue
        except Exception as e:
            print(f"Failed to parse output:\n{result.stdout}\nError: {e}")
            return None

    return total_time / num_runs

print("\n=== Sequential Runs ===")
seq_time = get_average_runtime("sequential")
print(f"Sequential average time: {seq_time:.4f} seconds\n")

bsp_speedups = []
for t in num_threads:
    print(f"\n=== BSP with {t} threads ===")
    avg_bsp_time = get_average_runtime("bsp", threads=t)
    speedup = seq_time / avg_bsp_time
    bsp_speedups.append(speedup)

worksteal_speedups = {rpt: [] for rpt in rows_per_task_options}
for rpt in rows_per_task_options:
    for t in num_threads:
        print(f"\n=== WorkSteal with {t} threads and {rpt} rows per task ===")
        avg_ws_time = get_average_runtime("worksteal", threads=t, rows_per_task=rpt)
        speedup = seq_time / avg_ws_time
        worksteal_speedups[rpt].append(speedup)

plt.figure(figsize=(10, 6))

plt.plot(num_threads, bsp_speedups, marker='o', label='BSP')

for rpt, speedups in worksteal_speedups.items():
    plt.plot(num_threads, speedups, marker='o', label=f'WorkSteal (rows per thread={rpt})')

plt.xlabel("Threads")
plt.ylabel("Speedup (Sequential Time / Parallel Time)")
plt.title("Speedup vs Threads")
plt.legend()
plt.grid(True)
plt.savefig("speedup_plot.png")