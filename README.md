# Dataset data source for Grafana

The Dataset Grafana data source plugin allows you to create and visualize graphs
and dashboards in Grafana using data in Dataset. You may want to use this plugin
to allow you to visualize Dataset data next to other data sources, for instance
when you want to monitor many feeds on a single dashboard.

![SystemDashboard](images/SystemDashboard.png)

With the Dataset plugin, you will be able to create and visualize your log-based
metrics along side all of your other data sources. It's a great way to have a
single pane of glass for today's complex systems. You can leverage Grafana alerts
based on Dataset data to notify you when there are possible issues. More
importantly, you'll soon be able to jump to Dataset's fast, easy and intuitive
platform to quickly identify the underlying causes of issues that may arise.

## Prerequisites

* **An installed Grafana server instance with write access**: This document
assumes that an existing instance of Grafana already exists. If you need help
bringing up a Grafana instance, please refer to the [documentation provided by
Grafana](https://grafana.com/docs/installation/).
* **A Dataset read log API Key**: A Dataset API key is required for Grafana to pull
data from Dataset. You can obtain one by going to your account in the Dataset
product and selecting the “API Keys” from the menu in the top right corner. You
can find documentation on API Keys [here](https://www.scalyr.com/help/api#scalyr-api-keys).

## Getting started

### Installing with grafana-cli

1. To install the stable version of the plugin using grafana-cli, run the following command:

   ```bash
   grafana-cli --pluginUrl \
   https://github.com/scalyr/scalyr-grafana-datasource-plugin/releases/download/3.0.0/sentinelone-dataset-datasource.zip \
   plugins install sentinelone-dataset-datasource
   ```

2. Update your Grafana configuration in the `grafana.ini` file to allow this plugin by adding the following line:

   ```bash
   allow_loading_unsigned_plugins = sentinelone-dataset-datasource
   ```

3. Adding plugins requires a restart of your grafana server.

    For init.d based services you can use the command:

    ```bash
    sudo service grafana-server restart
    ```

    For systemd based services you can use the following:

    ```bash
    systemctl restart grafana-server
    ```

If you require the development version, use the manual installation instructions.
### Installing manually

1. If you want a stable version of plugin, download the desired version from
[github releases](https://github.com/scalyr/scalyr-grafana-datasource-plugin/releases).
If you want the `development` version of the plugin,
clone the [plugin repository](https://github.com/scalyr/scalyr-grafana-datasource)
from GitHub. Switch to branch `go-rewrite-v2`

    ```bash
    git  clone https://github.com/scalyr/scalyr-grafana-datasource-plugin.git
    ```

2. Grafana plugins exist in the directory: `/var/lib/grafana/plugins/`. Create a folder for the dataset plugin:

    ```bash
    mkdir /var/lib/grafana/plugins/dataset
    ```

3. Copy the contents of the Dataset plugin into grafana:

    Stable version:

    ```bash
    tar -xvf scalyr_grafana_plugin_51057f6.tar.gz
    cp -rf dist/ /var/lib/grafana/plugins/scalyr/
    ```

    Development version:

    ```bash
    cp -r scalyr-grafana-datasource/dist/ /var/lib/grafana/plugins/scalyr/
    ```

4. Adding plugins requires a restart of your grafana server.

    For init.d based services you can use the command:

    ```bash
    sudo service grafana-server restart
    ```

    For systemd based services you can use the following:

    ```bash
    systemctl restart grafana-server
    ```
### Verify the Plugin was Installed

1. In order to verify proper installation you must log in to your grafana instance
   and navigate to **Configuration Settings -> Data Sources**.

    ![FirstImage](images/ConfigDataSource.png)

2. This will take you into the configuration page. If you already have other data
   sources installed, you will see them show up here. Click on the **Add data source** button:

    ![SecondImage](images/DataSoureConfig.png)

3. If you enter "Scalyr" in the search bar on the resulting page you should see “Scalyr Grafana
   Datasource” show up as an option.

    ![otherPlugin](images/SearchForPlugin.png)

4. Click on ***“Select”***. This will take you to a configuration page where you
   insert your API key mentioned in the prerequisite section.

    ![PluginConfig](images/PluginConfig.png)

5. Enter these settings:

    |Field Name | Value|
    | --- | --- |
    |Scalyr API Key | Your Scalyr Read Logs API Key|
    |Scalyr URL | `https://www.scalyr.com` or `https://eu.scalyr.com` for EU users.|

6. Click ***Save & Test*** to verify these settings are correct.