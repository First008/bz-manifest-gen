package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	//"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func introToPromptOptions() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome, here the options for you\n")
	fmt.Println("\t1: Create a Go Ethereum manifests for kubernetes\n")
	fmt.Println("\t2: Deploy generated manifests\n")
	fmt.Println("\n\t0: Quit\n")

	fmt.Println("So, what do you want to do? \n")

	opt, _ := getInput("", reader)

	switch opt {

	case "1":
		ethOptions()
	case "2":
		opt2, _ := getInput("\nWhich folder should i check for ymls? (ENTER default)", reader)
		if len(replace([]byte(opt2), "\n", "", -1)) == 0 {
			deployYmls("./eth-manifests")
		}

		deployYmls(opt2)
	case "0":
		os.Exit(0)

	default:
		fmt.Println("That was not a valid option...")
		introToPromptOptions()
	}
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		*files = append(*files, path)
		return nil
	}
}

func deployYmls(dir string) {

	config, clientset := getConfig()

	var files []string

	root := dir
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	// var statefulsets []string
	// var namespaceOfStatefulset string

	for _, file := range files[1:] {

		fmt.Println(file)

		yaml := readFile(file)

		// Create a YAML serializer.  JSON is a subset of YAML, so is supported too.
		s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

		// Decode the YAML to an object.
		var list corev1.List
		object, _, err := s.Decode(yaml, nil, &list)
		if err != nil {
			panic(err)
		}

		// Some types, e.g. List, contain RawExtensions.  If the appropriate types
		// are registered, these can be decoded in a second pass.
		for i, o := range list.Items {
			o.Object, _, err = s.Decode(o.Raw, nil, nil)
			if err != nil {
				panic(err)
			}
			o.Raw = nil

			list.Items[i] = o
		}

		//fmt.Printf("%#v\n", list)

		namespace, kind, _ := createObject(clientset, *config, object)
		// namespaceOfStatefulset = namespace

		if kind == "Job" {
			data := readFile("./eth-manifests/600.txt")
			init, _ := strconv.Atoi(string(data))
			time.Sleep(time.Duration(init*21) * time.Second)
		}

		// if kind == "StatefulSet" {
		// 	statefulsets = append(statefulsets, name)
		// }

		watch, err := clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal(err.Error())
		}

		go func() {
			for event := range watch.ResultChan() {
				p, ok := event.Object.(*v1.Pod)
				if !ok {
					log.Fatal("unexpected type")
				}
				//fmt.Println(p.Status.Phase)
				if string(p.Status.Phase) != "Running" {
					time.Sleep(10 * time.Second)
				}

			}
		}()

	}

	// for _, sts := range statefulsets {
	// 	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().String())
	// 	clientset.AppsV1().StatefulSets(namespaceOfStatefulset).Patch(context.Background(), sts, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})

	// }
}

func ethOptions() {
	reader := bufio.NewReader(os.Stdin)

	name, _ := getInput("\nPlease enter the name of the Blockchain: ", reader)
	nodeCountWanted, _ := getInput("\nPlease enter the amount of nodes you want: ", reader)
	storage, _ := getInput("\nPlease enter the storage space you want to reserve (GiB): ", reader)
	nameTX, _ := getInput("\nPlease enter the name of the TX node: ", reader)
	id, _ := getInput("\nPlease enter the id of blockchain: ", reader)
	pwd, _ := getInput("\nPlease set the accounts passwd: ", reader)

	g := newGeth(name)

	var nodeCount int

	nodeCount, _ = strconv.Atoi(nodeCountWanted)

	g.setNodeCount(nodeCount)
	g.setStorageCapacity(storage)
	g.setNameOfTX(nameTX)
	g.setid(id)
	g.setAccountPwd(pwd)

	prepare000(g)
	prepare010(g)
	prepare100(g)
	prepare110(g)
	prepare200(g)
	prepare210(g)
	prepare220(g)
	prepare300(g)
	prepare310(g)
	prepare320(g)
	prepare330(g)
	prepare340(g)
	prepare500(g)

	introToPromptOptions()
}
