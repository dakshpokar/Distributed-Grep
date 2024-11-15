import matplotlib.pyplot as plt
import numpy as np

queries = [1, 2, 3, 4]
query_1 = [500.048882, 498.139625, 520.314082, 501.736553, 501.609901]
query_2 = [297.582278, 298.411494, 298.877676, 298.529602, 284.993645]
query_3 = [704.443657, 688.471485, 703.70245, 701.9439, 699.762611]
query_4 = [138.387803, 137.257963, 139.113749, 136.053214, 136.975515]
avg_latency = [np.mean(query_1), np.mean(query_2), np.mean(query_3), np.mean(query_4)]
std_dev = [np.std(query_1), np.std(query_2), np.std(query_3), np.std(query_4)]
fig, ax = plt.subplots()

bar = ax.bar(queries, avg_latency, yerr=std_dev, capsize=5, alpha=0.7)

ax.set_xlabel("queries")
ax.set_ylabel("Average Query Latency (ms)")
ax.set_title("Query Latency for 4 Machines (60 MB log files each)")

ax.set_xticks(queries)
ax.set_xticklabels(["PUT", "PUT.*300", "PUT|DELETE", "INFO"])

for i, bar in enumerate(bar):
    height = bar.get_height()
    ax.text(bar.get_x() + bar.get_width()/2, height/1.5,
            f'Mean: {avg_latency[i]:.1f}\n\nSD: {std_dev[i]:.1f}',
            ha='center', va='bottom')

ax.grid(True, linestyle="--", alpha=0.7)

plt.tight_layout()
plt.show()
