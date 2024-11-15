import subprocess

REMOTE_USER = "root"
REMOTE_DIR = "/root/mp1-distributed-systems"


# List of remote hosts
REMOTE_HOSTS = [
    "fa24-cs425-9201.cs.illinois.edu",
    "fa24-cs425-9202.cs.illinois.edu",
    "fa24-cs425-9203.cs.illinois.edu",
    "fa24-cs425-9204.cs.illinois.edu",
    "fa24-cs425-9205.cs.illinois.edu",
    "fa24-cs425-9206.cs.illinois.edu",
    "fa24-cs425-9207.cs.illinois.edu",
    "fa24-cs425-9208.cs.illinois.edu",
    "fa24-cs425-9209.cs.illinois.edu",
    "fa24-cs425-9210.cs.illinois.edu",
]

def generate_log_files_on_all_machines():
    for index, host in enumerate(REMOTE_HOSTS):
        print(f"Generating log files on {host}")
        subprocess.run(
            [
                "ssh",
                f"{REMOTE_USER}@{host}",
                f"cd {REMOTE_DIR} &&",
                f"rm -rf test_log{index+1}.log &&",
                f"nohup go run unit_tests/scripts/generate_log_files_on_machine.go {index+1} & exit",
            ]
        )
        print(f"Log files generated on {host}")
    print("Log file generation process completed...")

def run_distributed_grep(pattern):
    result = subprocess.run(
            [
                "ssh",
                f"{REMOTE_USER}@fa24-cs425-9201.cs.illinois.edu",
                f"cd {REMOTE_DIR} &&",
                f"go run client.go {pattern} {REMOTE_DIR}/test_log",
            ],
            capture_output=True,
            text=True,
        )
    # Save result.stdout to a file
    with open("output_" + pattern + ".txt", "w") as f:
        f.write(result.stdout)
    

    output = result.stdout.split("###########################")[1]

    filtered_output = output.split("\n")
    m = {}
    for line in filtered_output:
        if line:
            key, value = line.split(" ::: ")
            m[key.strip()] = value.strip()
    return m

def generate_test_logs(result, test_case_map, pattern_type):
    passed = 0
    for key, value in result.items():
        stripped_key = key.split(":")[1].strip()
        expected_value = test_case_map[stripped_key + ".cs.illinois.edu"]

        if "dead" not in value:
            current_value = int(value)
            if current_value == expected_value:
                passed += 1
                print(f"\t✅ Test case passed for {key}", current_value, expected_value)
            else:
                print(f"\t❌ Test case failed for {key}", current_value, expected_value)
        else:
            print(f"Test case failed for {key}")

    if passed == 10:
        print("✅ All test cases passed successfully")
    else:
        print(f"❌ Only {passed} test cases passed successfully for {pattern_type} Pattern")

def test_frequent_pattern():
    test_case_1_frequent_pattern = {
        "fa24-cs425-9201.cs.illinois.edu": 56879,
        "fa24-cs425-9202.cs.illinois.edu": 53604,
        "fa24-cs425-9203.cs.illinois.edu": 54063,
        "fa24-cs425-9204.cs.illinois.edu": 53921,
        "fa24-cs425-9205.cs.illinois.edu": 53935,
        "fa24-cs425-9206.cs.illinois.edu": 53820,
        "fa24-cs425-9207.cs.illinois.edu": 53538,
        "fa24-cs425-9208.cs.illinois.edu": 54880,
        "fa24-cs425-9209.cs.illinois.edu": 54029,
        "fa24-cs425-9210.cs.illinois.edu": 53435,
    }
    print("Testing frequent pattern - PUT ...")
    
    result = run_distributed_grep("PUT")
    
    generate_test_logs(result, test_case_1_frequent_pattern, "Frequent")
    
def test_rare_pattern():
    # PUT.*300
    test_case_2_rate_pattern = {
        "fa24-cs425-9201.cs.illinois.edu": 15,
        "fa24-cs425-9202.cs.illinois.edu": 19,
        "fa24-cs425-9203.cs.illinois.edu": 21,
        "fa24-cs425-9204.cs.illinois.edu": 15,
        "fa24-cs425-9205.cs.illinois.edu": 17,
        "fa24-cs425-9206.cs.illinois.edu": 11,
        "fa24-cs425-9207.cs.illinois.edu": 11,
        "fa24-cs425-9208.cs.illinois.edu": 9,
        "fa24-cs425-9209.cs.illinois.edu": 11,
        "fa24-cs425-9210.cs.illinois.edu": 15,
    }

    print("Testing rare pattern - PUT.*300 ...")
    
    result = run_distributed_grep("PUT.*300")
    
    generate_test_logs(result, test_case_2_rate_pattern, "Rare")

def test_somewhat_frequent_pattern():
    # PUT.*301
    test_case_3_somewhat_frequent_pattern = {
        "fa24-cs425-9201.cs.illinois.edu": 2255,
        "fa24-cs425-9202.cs.illinois.edu": 2159,
        "fa24-cs425-9203.cs.illinois.edu": 2255,
        "fa24-cs425-9204.cs.illinois.edu": 2162,
        "fa24-cs425-9205.cs.illinois.edu": 2118,
        "fa24-cs425-9206.cs.illinois.edu": 2248,
        "fa24-cs425-9207.cs.illinois.edu": 2214,
        "fa24-cs425-9208.cs.illinois.edu": 2261,
        "fa24-cs425-9209.cs.illinois.edu": 2177,
        "fa24-cs425-9210.cs.illinois.edu": 2185,
    }
    print("Testing Somewhat Frequent Pattern - PUT.*301 ...")
    
    result = run_distributed_grep("PUT.*301")
    
    generate_test_logs(result, test_case_3_somewhat_frequent_pattern, "Somewhat Frequent")

def test_specific_line():
    # PUT.*700
    test_case_4_specific_line = {
        "fa24-cs425-9201.cs.illinois.edu": 0,  
        "fa24-cs425-9202.cs.illinois.edu": 1,  
        "fa24-cs425-9203.cs.illinois.edu": 0,  
        "fa24-cs425-9204.cs.illinois.edu": 0, 
        "fa24-cs425-9205.cs.illinois.edu": 0, 
        "fa24-cs425-9206.cs.illinois.edu": 0, 
        "fa24-cs425-9207.cs.illinois.edu": 0,  
        "fa24-cs425-9208.cs.illinois.edu": 0,  
        "fa24-cs425-9209.cs.illinois.edu": 0,  
        "fa24-cs425-9210.cs.illinois.edu": 0,  
    }
    print("Testing Specific Line - GET.*700 ...")
    result = run_distributed_grep("GET.*700")
    
    generate_test_logs(result, test_case_4_specific_line, "Specific Line")

if __name__ == "__main__":
    generate_log_files_on_all_machines()
    test_frequent_pattern()
    test_rare_pattern()
    test_somewhat_frequent_pattern()
    test_specific_line()
