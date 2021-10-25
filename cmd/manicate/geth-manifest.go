package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	appsv1 "github.com/openshift/api/apps/v1"
	authorizationv1 "github.com/openshift/api/authorization/v1"
	buildv1 "github.com/openshift/api/build/v1"
	imagev1 "github.com/openshift/api/image/v1"
	networkv1 "github.com/openshift/api/network/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	projectv1 "github.com/openshift/api/project/v1"
	quotav1 "github.com/openshift/api/quota/v1"
	routev1 "github.com/openshift/api/route/v1"
	securityv1 "github.com/openshift/api/security/v1"
	templatev1 "github.com/openshift/api/template/v1"
	userv1 "github.com/openshift/api/user/v1"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func getConfig() (*rest.Config, *kubernetes.Clientset) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})

	config, err := kubeconfig.ClientConfig()

	if err != nil {
		panic(err)
	}

	clierntset := kubernetes.NewForConfigOrDie(config)

	return config, clierntset
}

func init() {
	// The Kubernetes Go client (nested within the OpenShift Go client)
	// automatically registers its types in scheme.Scheme, however the
	// additional OpenShift types must be registered manually.  AddToScheme
	// registers the API group types (e.g. route.openshift.io/v1, Route) only.
	appsv1.AddToScheme(scheme.Scheme)
	authorizationv1.AddToScheme(scheme.Scheme)
	buildv1.AddToScheme(scheme.Scheme)
	imagev1.AddToScheme(scheme.Scheme)
	networkv1.AddToScheme(scheme.Scheme)
	oauthv1.AddToScheme(scheme.Scheme)
	projectv1.AddToScheme(scheme.Scheme)
	quotav1.AddToScheme(scheme.Scheme)
	routev1.AddToScheme(scheme.Scheme)
	securityv1.AddToScheme(scheme.Scheme)
	templatev1.AddToScheme(scheme.Scheme)
	userv1.AddToScheme(scheme.Scheme)

	// If you need to serialize/deserialize legacy (non-API group) OpenShift
	// types (e.g. v1, Route), these must be additionally registered using
	// AddToSchemeInCoreGroup.
	appsv1.AddToSchemeInCoreGroup(scheme.Scheme)
	authorizationv1.AddToSchemeInCoreGroup(scheme.Scheme)
	buildv1.AddToSchemeInCoreGroup(scheme.Scheme)
	imagev1.AddToSchemeInCoreGroup(scheme.Scheme)
	networkv1.AddToSchemeInCoreGroup(scheme.Scheme)
	oauthv1.AddToSchemeInCoreGroup(scheme.Scheme)
	projectv1.AddToSchemeInCoreGroup(scheme.Scheme)
	quotav1.AddToSchemeInCoreGroup(scheme.Scheme)
	routev1.AddToSchemeInCoreGroup(scheme.Scheme)
	securityv1.AddToSchemeInCoreGroup(scheme.Scheme)
	templatev1.AddToSchemeInCoreGroup(scheme.Scheme)
	userv1.AddToSchemeInCoreGroup(scheme.Scheme)

}

func createObject(kubeClientset kubernetes.Interface, restConfig rest.Config, obj runtime.Object) (string, string, string) {
	// Create a REST mapper that tracks information about the available resources in the cluster.
	groupResources, err := restmapper.GetAPIGroupResources(kubeClientset.Discovery())
	if err != nil {
		panic(err)
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	// Get some metadata needed to make the REST request.
	gvk := obj.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		panic(err)
	}

	namespace, _ := meta.NewAccessor().Namespace(obj)
	name, _ := meta.NewAccessor().Name(obj)
	kind, err := meta.NewAccessor().Kind(obj)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n[STATUS]: Applying %q under namespace %q\t\t From:", name, namespace)

	// Create a client specifically for creating the object.
	restClient, err := newRestClient(restConfig, mapping.GroupVersionKind.GroupVersion())
	if err != nil {
		panic(err)
	}

	// Use the REST helper to create the object in the "default" namespace.
	restHelper := resource.NewHelper(restClient, mapping)
	restHelper.Create(namespace, false, obj)
	return namespace, kind, name
}

func newRestClient(restConfig rest.Config, gv schema.GroupVersion) (rest.Interface, error) {
	restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	restConfig.GroupVersion = &gv
	if len(gv.Group) == 0 {
		restConfig.APIPath = "/api"
	} else {
		restConfig.APIPath = "/apis"
	}

	return rest.RESTClientFor(&restConfig)
}

func moveToOld() {

	if _, err := os.Stat("./eth-manifests"); !os.IsNotExist(err) {
		if _, err := os.Stat("./eth-manifests-old"); !os.IsNotExist(err) {
			os.RemoveAll("./eth-manifests-old/")
		}
		os.Rename("./eth-manifests", "./eth-manifests-old")
	}

}

func readFile(s string) []byte {

	data, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Println("File reading error", err)
		os.Exit(1)
	}

	return data
}

func writeFile(path string, data []byte) {

	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}

}

func replace(data []byte, old string, new string, n int) []byte {

	return []byte(strings.Replace(string(data), old, new, n))

}

func prepare000(g geth) {

	moveToOld() // Moves old manifests to old named dir

	err := os.Mkdir("./eth-manifests", 0755)
	if err != nil {
		log.Fatal(err)
	}

	data := readFile("./GAS/000-geth-namespace.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	writeFile("./eth-manifests/000.yml", data)
}

func prepare010(g geth) {

	data := readFile("./GAS/010-geth-storage.yml")

	//data = []byte(strings.Replace(string(data), "$(NAMESPACE)", g.name, 2))

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	writeFile("./eth-manifests/010.yml", data)

}

func prepare100(g geth) {

	data := readFile("./GAS/100-geth-pv.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	data = replace(data, "$(STORAGE)", g.storage, -1)

	writeFile("./eth-manifests/100.yml", data)
}

func prepare110(g geth) {

	data := readFile("./GAS/110-geth-pvc.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	data = replace(data, "$(STORAGE)", g.storage, -1)

	writeFile("./eth-manifests/110.yml", data)

}

func prepare200(g geth) {

	data := readFile("./GAS/200-etherstats-service.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	writeFile("./eth-manifests/200.yml", data)

}

func prepare210(g geth) {

	data := readFile("./GAS/210-etherstats-secret.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	writeFile("./eth-manifests/210.yml", data)

}

func prepare220(g geth) {

	data := readFile("./GAS/220-etherstats-dashb.yaml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	writeFile("./eth-manifests/220.yml", data)

}
func prepare300(g geth) {

	data := readFile("./GAS/300-geth-confmap.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	data = replace(data, "$(PWD)", g.accountPwd, -1)

	var nodesToConfigMap strings.Builder

	for index := 1; index <= g.nodeCount; index++ {

		nodesToConfigMap.WriteString(fmt.Sprintf("node%d: \"node-%d\"\n  ", index, index))

	}

	var addressesAndBalances strings.Builder

	for index := 1; index <= g.nodeCount; index++ {

		addressesAndBalances.WriteString(fmt.Sprintf(
			`        "addres%ds": {%s          "balance": "0x200000000000000000000000000000000000000000000000000000000000000"%s        },`, index, "\n", "\n"))
		addressesAndBalances.WriteString("\n")

	}

	data = replace(data, "$(ADDRESSES)", strings.TrimRight(strings.TrimRight(addressesAndBalances.String(), "\n"), ","), -1)

	data = replace(data, "$(ID)", g.id, -1)

	data = replace(data, "$(NODES): $(NODES)", nodesToConfigMap.String(), -1)

	writeFile("./eth-manifests/300.yml", data)

}

func prepare310(g geth) {

	data := readFile("./GAS/310-account-secret.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	writeFile("./eth-manifests/310.yml", data)
}

func prepare320(g geth) {

	data := readFile("./GAS/320-geth-service.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	var portsFromNodes strings.Builder

	tcpPort := 8545
	udpPort := 30311

	for index := 1; index <= g.nodeCount; index++ {

		portsFromNodes.WriteString(fmt.Sprintf("- port: %d\n    name: node%d\n  - port: %d\n    protocol: UDP\n    name: udp-node-%d\n  ", tcpPort, index, udpPort, index))
		tcpPort += 1
		udpPort += 1
	}

	data = replace(data, "$(PORTS): $(NODES)", portsFromNodes.String(), -1)

	writeFile("./eth-manifests/320.yml", data)

}

func prepare330(g geth) {

	data := readFile("./GAS/330-geth-job.yml")

	var envFromNodes strings.Builder

	for index := 1; index <= g.nodeCount; index++ {

		envFromNodes.WriteString(fmt.Sprintf("- name: node%d\n          valueFrom:\n            configMapKeyRef:\n              name: %s-geth\n              key: node%d\n        ", index, g.name, index))
	}

	var envFromNodesN strings.Builder

	for index := 1; index <= g.nodeCount; index++ {

		envFromNodesN.WriteString(fmt.Sprintf("- name: node%d\n            valueFrom:\n              configMapKeyRef:\n                name: %s-geth\n                key: node%d\n          ", index, g.name, index))
	}

	createAccountsForNodes :=
		`      - name: create-account-node-$(NUMBER)
        image: ethereum/client-go:release-1.8
        imagePullPolicy: IfNotPresent
        command: ["/bin/sh"]
        args:
        - "-c"
        - "printf '$(ACCOUNT_SECRET)\n$(ACCOUNT_SECRET)\n' | geth account new --datadir /root/.ethereum/$(NAMESPACE)/$(node$(NUMBER))"
        env:
        - name: ACCOUNT_SECRET
          valueFrom:
            configMapKeyRef:
              name: $(NAMESPACE)-geth
              key: password.txt
        - name: node$(NUMBER)
          valueFrom:
            configMapKeyRef:
              name: $(NAMESPACE)-geth
              key: node$(NUMBER)

        volumeMounts:
          - name: data2
            mountPath: /root/.ethereum/$(NAMESPACE)`

	initGenesisConteinerString :=
		`      - name: init-genesis-node$(NUMBER)
        image: ethereum/client-go:release-1.8
        imagePullPolicy: IfNotPresent
        args: [
          "--datadir", "/root/.ethereum/$(NAMESPACE)/$(node$(NUMBER))",
          "init", "/root/.ethereum/$(NAMESPACE)/configuredgenesis.json"
          ]

        env:
        - name: node$(NUMBER)
          valueFrom:
            configMapKeyRef:
              name: $(NAMESPACE)-geth
              key: node$(NUMBER)

        volumeMounts:
        - name: data2
          mountPath: /root/.ethereum/$(NAMESPACE)
        - name: config-genesis
          mountPath: /root/.ethereum/$(NAMESPACE)/genesis.json
          subPath: genesis.json
        - name: config-genesis
          mountPath: /root/.ethereum/$(NAMESPACE)/pswd/password.txt
          subPath: password.txt`

	var initGenesisConteiners strings.Builder

	for index := 1; index <= g.nodeCount; index++ {
		init := strings.ReplaceAll(initGenesisConteinerString, "$(NUMBER)", fmt.Sprintf("%d", index))
		initGenesisConteiners.WriteString(init)
		initGenesisConteiners.WriteString("\n\n")
	}

	var createAccountString strings.Builder

	for index := 1; index <= g.nodeCount; index++ {
		init := strings.ReplaceAll(createAccountsForNodes, "$(NUMBER)", fmt.Sprintf("%d", index))
		createAccountString.WriteString(init)
		createAccountString.WriteString("\n\n")
	}

	var getAccountAdresses0x strings.Builder

	for index := 1; index <= g.nodeCount; index++ {
		getAccountAdresses0x.WriteString(fmt.Sprintf("printf '$(ACCOUNT_SECRET)*n' | echo -n `ethkey inspect /root/.ethereum/$(NAMESPACE)/$(node%d)/keystore/*` > /root/.ethereum/$(NAMESPACE)/$(node%d)/keys.txt && grep -Eo '0x.{40}' /root/.ethereum/$(NAMESPACE)/$(node%d)/keys.txt | head -1  > /root/.ethereum/$(NAMESPACE)/$(node%d)/address0x.txt &&", index, index, index, index))
	}

	var getAccountAdresses strings.Builder

	for index := 1; index <= g.nodeCount; index++ {
		getAccountAdresses.WriteString(fmt.Sprintf("cut -c 3- /root/.ethereum/$(NAMESPACE)/$(node%d)/address0x.txt > /root/.ethereum/$(NAMESPACE)/$(node%d)/address.txt &&", index, index))
	}

	var genesisConfiguration strings.Builder

	for index := 1; index <= g.nodeCount; index++ {
		genesisConfiguration.WriteString(fmt.Sprintf("echo -n `cat /root/.ethereum/$(NAMESPACE)/$(node%d)/address.txt` >> /root/.ethereum/$(NAMESPACE)/newaddresses.txt && ", index))
	}

	genesisConfiguration.WriteString("cp /root/.ethereum/$(NAMESPACE)/genesis.json /root/.ethereum/$(NAMESPACE)/configuredgenesis.json && sed -i -e  s/newaddresses/`cat /root/.ethereum/$(NAMESPACE)/newaddresses.txt`/g /root/.ethereum/$(NAMESPACE)/configuredgenesis.json && ")

	index := 1
	for index <= g.nodeCount {
		genesisConfiguration.WriteString(fmt.Sprintf("sed -i -e  s/addres%ds/`cat /root/.ethereum/$(NAMESPACE)/$(node%d)/address.txt`/g /root/.ethereum/$(NAMESPACE)/configuredgenesis.json && ", index, index))
		index++
	}

	var enodeFile strings.Builder
	index = 1
	udpPort := 30311
	for index <= g.nodeCount {
		enodeFile.WriteString(fmt.Sprintf("bootnode --genkey /root/.ethereum/$(NAMESPACE)/$(node%d)/bootnode.key && echo -n enode://`bootnode --nodekey /root/.ethereum/$(NAMESPACE)/$(node%d)/bootnode.key --writeaddress`@10.152.183.132:%d, >> /root/.ethereum/$(NAMESPACE)/enodelist.txt && ", index, index, udpPort))
		index++
		udpPort++
	}

	enodeFile.WriteString("sed -i -e '$ s/,$//' /root/.ethereum/$(NAMESPACE)/enodelist.txt")

	// bootnode --genkey /root/.ethereum/$(NAMESPACE)/$(node1)/bootnode.key && echo -n enode://`bootnode --nodekey /root/.ethereum/$(NAMESPACE)/$(node1)/bootnode.key --writeaddress`@10.152.183.132:30311, > /root/.ethereum/$(NAMESPACE)/enodelist.txt &&
	// $(ENODEARG)
	data = replace(data, "$(ENV)", envFromNodes.String(), -1)
	data = replace(data, "$(NENV)", envFromNodesN.String(), -1)
	data = replace(data, "$(ENODEARG)", enodeFile.String(), -1)
	data = replace(data, "$(CA)", strings.TrimRight(strings.TrimRight(createAccountString.String(), " "), "&"), -1)
	data = replace(data, "$(GAA0x)", strings.TrimRight(strings.TrimRight(getAccountAdresses0x.String(), " "), "&"), -1)
	data = replace(data, "$(GAA)", strings.TrimRight(strings.TrimRight(getAccountAdresses.String(), " "), "&"), -1)
	data = replace(data, "$(INITGENESISCONTEINERS)", initGenesisConteiners.String(), -1)
	data = replace(data, "$(GC)", strings.TrimRight(strings.TrimRight(genesisConfiguration.String(), " "), "&"), -1)

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	writeFile("./eth-manifests/330.yml", replace(data, "*n", `\n`, -1))
}

func prepare340(g geth) {

	var conPort int = 8545
	var tcpPort int = 30311
	var udpPort int = 30311
	index := 1

	for index <= g.nodeCount {
		var data []byte
		var nodeStatefulSet strings.Builder

		data = readFile("./GAS/340-geth-node.yml")
		data = replace(data, "$(NUMBER)", fmt.Sprintf("%d", index), -1)
		//data = replace(data, "$(ARG)", "geth --bootnodes $(cat /root/.ethereum/$(NAMESPACE)/enodelist.txt) --datadir /root/.ethereum/$(NAMESPACE)/$(NODE_NAME) --syncmode full --port $(UDPPORT) --http --http.addr 0.0.0.0 --http.port $(CONPORT) --http.vhosts '*' --http.api admin,eth,miner,net,txpool,personal,web3,clique --http.corsdomain '*' --networkid $(NETWORK_ID) --ethstats=$(HOSTNAME):$(ETHSTATS_SECRET)@$(ETHSTATS_SVC) --allow-insecure-unlock --unlock `cat /root/.ethereum/$(NAMESPACE)/$(NODE_NAME)/address0x.txt` --password /root/.ethereum/$(NAMESPACE)/pswd/password.txt --miner.gasprice 0 --miner.etherbase `cat /root/.ethereum/$(NAMESPACE)/$(NODE_NAME)/address0x.txt` --mine", -1)
		data = replace(data, "$(CONPORT)", fmt.Sprintf("%d", conPort), -1)
		data = replace(data, "$(TCPPORT)", fmt.Sprintf("%d", tcpPort), -1)
		data = replace(data, "$(UDPPORT)", fmt.Sprintf("%d", udpPort), -1)
		nodeStatefulSet.WriteString(string(data))
		nodeStatefulSet.WriteString("\n")

		tcpPort = tcpPort + 1
		udpPort = udpPort + 1
		conPort = conPort + 1

		data = replace([]byte(nodeStatefulSet.String()), "$(NAMESPACE)", g.name, -1)

		writeFile(fmt.Sprintf("./eth-manifests/340-%d.yml", index), data)
		index = index + 1

	}
	writeFile("./eth-manifests/600.txt", []byte(strconv.Itoa(index)))

}

func prepare500(g geth) {
	data := readFile("./GAS/500-alpine-for-monitoring.yml")

	data = replace(data, "$(NAMESPACE)", g.name, -1)

	writeFile("./eth-manifests/500.yml", data)

}
